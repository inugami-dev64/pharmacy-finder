package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pharmafinder/types"
	"pharmafinder/utils"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
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

type EmptyBody struct{}

// HttpRequestDetails is a struct that contains relevant data about
// the request that was made
type HttpRequestDetails[B interface{}] struct {
	Path   string
	Method string

	// Unmarshalled HTTP request body
	Body B

	// URL parameters
	Params url.Values

	// URL path variables
	// e.g. /api/v1/{var1}/{var2}
	PathVars map[string]string
}

type CallbackFunction[T interface{}, B interface{}] = func(details *HttpRequestDetails[B]) (int, interface{}, error)

type HttpRequestHandler[T interface{}, B interface{}] struct {
	callback CallbackFunction[T, B]
	pattern  string
	validate *validator.Validate
	methods  []string
	logger   zerolog.Logger
}

func NewRequestsHandler[T interface{}, B interface{}](
	callback CallbackFunction[T, B],
	pattern string,
	methods []string) Route {
	return &HttpRequestHandler[T, B]{
		callback: callback,
		pattern:  pattern,
		validate: validator.New(),
		methods:  methods,
		logger:   utils.GetLogger("WEB"),
	}
}

func (handler *HttpRequestHandler[T, B]) Pattern() string {
	return handler.pattern
}

func (handler *HttpRequestHandler[T, B]) Methods() []string {
	return handler.methods
}

func (handler *HttpRequestHandler[T, B]) assignBody(r *http.Request, w http.ResponseWriter, details *HttpRequestDetails[B]) {
	// check if given request body should be unmarshalled
	if _, ok := any(details.Body).(EmptyBody); !ok {
		badRequestLogEvent := handler.logger.Warn().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("addr", r.RemoteAddr).
			Int("code", http.StatusBadRequest)

		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			badRequestLogEvent.Msgf("Expected content-type is application/json, got %s", r.Header.Get("Content-Type"))
			createJsonResponse(w, http.StatusBadRequest, types.HttpError{
				StatusCode: http.StatusBadRequest,
				Timestamp:  types.Time(time.Now().UTC()),
				Message:    fmt.Sprintf("Expected content-type is application/json, got %s", r.Header.Get("Content-Type")),
			})
			return
		}

		var body B
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			badRequestLogEvent.Msgf("Failed to read request body: %v", err)

			createJsonResponse(w, http.StatusBadRequest, types.HttpError{
				StatusCode: http.StatusBadRequest,
				Timestamp:  types.Time(time.Now().UTC()),
				Message:    "Invalid request body",
			})
			return
		}

		err = json.Unmarshal(bytes, &body)
		if err != nil {
			badRequestLogEvent.Msgf("Failed to json unmarshal request body: %v", err)

			createJsonResponse(w, http.StatusBadRequest, types.HttpError{
				StatusCode: http.StatusBadRequest,
				Timestamp:  types.Time(time.Now().UTC()),
				Message:    "Malformed JSON body",
			})
			return
		}

		details.Body = body
	}
}

func (handler *HttpRequestHandler[T, B]) validateBody(r *http.Request, w http.ResponseWriter, body *B) bool {
	if _, ok := any(body).(*EmptyBody); !ok {
		err := handler.validate.Struct(*body)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			handler.logger.Warn().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("addr", r.RemoteAddr).
				Int("code", http.StatusBadRequest).
				Msgf("Validation error: %s", errs[0].Error())

			createJsonResponse(w, http.StatusBadRequest, types.HttpError{
				StatusCode: http.StatusInternalServerError,
				Timestamp:  types.Time(time.Now().UTC()),
				Message:    errs[0].Error(),
			})
			return false
		}
	}

	return true
}

func (handler *HttpRequestHandler[T, B]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	details := HttpRequestDetails[B]{
		Path:     r.URL.Path,
		Method:   r.Method,
		Params:   r.URL.Query(),
		PathVars: mux.Vars(r),
	}

	// in cases where we have a request body provided, we perform json unmarshalling
	// and data validation
	handler.assignBody(r, w, &details)
	if !handler.validateBody(r, w, &details.Body) {
		return
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
	w.Header().Add("Content-Type", "application/json; utf-8")
	w.WriteHeader(code)
	w.Write(b)
}
