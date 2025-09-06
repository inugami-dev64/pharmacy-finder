package web

import (
	"encoding/json"
	"net/http"
	"net/url"
	"pharmafinder/types"
	"pharmafinder/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

// Controller interface defines the shared functionality
// for all controller structs
type Controller interface {
	GetRoutes() []Route
}

// Route is an http.Handler that knows the mux pattern
// under which it will be registered
type Route interface {
	http.Handler

	// Pattern reports the relative path at which this is registered
	Pattern() string

	// Methods reports all HTTP methods that this handler accepts
	Methods() []string
}

// HttpRequestDetails is a struct that contains relevant data about
// the request that was made
type HttpRequestDetails struct {
	Path   string
	Method string

	// URL parameters
	Params url.Values

	// URL path variables
	// e.g. /api/v1/{var1}/{var2}
	PathVars map[string]string
}

type CallbackFunction[T interface{}] = func(details *HttpRequestDetails) (int, interface{}, error)

type HttpRequestHandler[T interface{}] struct {
	callback CallbackFunction[T]
	pattern  string
	methods  []string
	logger   zerolog.Logger
}

func NewRequestsHandler[T interface{}](
	callback CallbackFunction[T],
	pattern string,
	methods []string) Route {
	return &HttpRequestHandler[T]{
		callback: callback,
		pattern:  pattern,
		methods:  methods,
		logger:   utils.GetLogger("WEB"),
	}
}

func (handler *HttpRequestHandler[T]) Pattern() string {
	return handler.pattern
}

func (handler *HttpRequestHandler[T]) Methods() []string {
	return handler.methods
}

func (handler *HttpRequestHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	details := HttpRequestDetails{
		Path:     r.URL.Path,
		Method:   r.Method,
		Params:   r.URL.Query(),
		PathVars: mux.Vars(r),
	}

	code, resp, err := handler.callback(&details)
	if err != nil {
		createJsonResponse(w, http.StatusInternalServerError, types.HttpError{
			StatusCode: http.StatusInternalServerError,
			Timestamp:  types.Time(time.Now().UTC()),
			Message:    "Internal server error",
		})
		handler.logger.Debug().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("addr", r.RemoteAddr).
			Int("code", http.StatusInternalServerError).
			Msg("Request made")
		return
	}

	createJsonResponse(w, code, resp)
	handler.logger.Debug().
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("addr", r.RemoteAddr).
		Int("code", code).
		Msg("Request made")
}

// Utility functions down below

func createJsonResponse(w http.ResponseWriter, code int, resp interface{}) {
	b, _ := json.Marshal(resp)
	w.WriteHeader(code)
	w.Write(b)
}
