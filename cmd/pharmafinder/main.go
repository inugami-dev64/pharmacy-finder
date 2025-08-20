package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"pharmafinder"
	v1 "pharmafinder/api/v1"
	"pharmafinder/db"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func NewServerMux(routes []Route) *mux.Router {
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	for _, route := range routes {
		apiRouter.Handle(route.Pattern(), route).Methods(route.Methods()...)
	}

	r.PathPrefix("/").
		Methods("GET").
		HandlerFunc(pharmafinder.StaticServer)

	return r
}

func NewHTTPServer(lc fx.Lifecycle, mux *mux.Router) *http.Server {
	server := &http.Server{
		Handler:      mux,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			log.Println("Starting HTTP server at", server.Addr)
			go server.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return server
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

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

func main() {
	// Attempt to load .env files if they exist
	godotenv.Load("deploy/.env")
	fx.New(
		fx.Provide(
			NewHTTPServer,
			fx.Annotate(
				NewServerMux,
				fx.ParamTags(`group:"routes"`),
			),

			// Data access layer
			db.ProvideDatabaseHandle,
			db.ProvidePharmacyRepository,

			// /pharmacies/* handlers
			AsRoute(v1.NewPharmaciesHandler),
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
