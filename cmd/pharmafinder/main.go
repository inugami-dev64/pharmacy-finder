package main

import (
	"context"
	"net"
	"net/http"
	"pharmafinder"
	v1 "pharmafinder/api/v1"
	"pharmafinder/bg"
	"pharmafinder/db"
	"pharmafinder/utils"
	"pharmafinder/web"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func NewServerMux(routes [][]web.Route) *mux.Router {
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	for _, routeGroup := range routes {
		for _, route := range routeGroup {
			apiRouter.Handle(route.Pattern(), route).Methods(route.Methods()...)
		}
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
			logger := utils.GetLogger("SRV")
			ln, err := net.Listen("tcp", server.Addr)
			if err != nil {
				return err
			}
			logger.Info().Msgf("Starting HTTP server at %s", server.Addr)
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
		fx.WithLogger(func() fxevent.Logger {
			return &utils.FXZerologLogger{Logger: utils.GetLogger("FX")}
		}),
		fx.Provide(
			NewHTTPServer,
			fx.Annotate(
				NewServerMux,
				fx.ParamTags(`group:"routes"`),
			),

			// Data access layer
			db.ProvideDatabaseHandle,
			db.ProvidePharmacyRepository,

			// Utilities
			utils.ProvideHTTPClient,

			// Background workers
			fx.Annotate(
				bg.ProvideBenuScraper,
				fx.ResultTags(`group:"scrapers"`),
			),
			fx.Annotate(
				bg.ProvideApothekaScraper,
				fx.ResultTags(`group:"scrapers"`),
			),
			fx.Annotate(
				bg.ProvideSydameapteekScraper,
				fx.ResultTags(`group:"scrapers"`),
			),
			fx.Annotate(
				bg.NewCronJob,
				fx.ParamTags(`group:"scrapers"`),
			),

			// /pharmacies/* controller
			fx.Annotate(
				v1.ProvidePharmacyController,
				fx.ResultTags(`group:"routes"`),
			),
		),
		fx.Invoke(func(*http.Server, bg.CronJob) {}),
	).Run()
}
