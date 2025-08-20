package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"pharmafinder"
	"pharmafinder/db"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func NewServerMux() *mux.Router {
	r := mux.NewRouter()
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

func main() {
	// Attempt to load .env files if they exist
	godotenv.Load("deploy/.env")
	fx.New(
		fx.Provide(NewHTTPServer),
		fx.Provide(NewServerMux),
		fx.Provide(db.ProvideDatabaseHandle),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
