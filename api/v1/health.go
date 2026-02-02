package v1

import (
	"net/http"
	"pharmafinder/db"
	"pharmafinder/db/dto"
	"pharmafinder/service"
	"pharmafinder/web"

	"github.com/jmoiron/sqlx"
)

type HealthCheckController struct {
	db        *sqlx.DB
	repo      db.ModeratorUserRepository
	collector service.ScraperStatCollector
}

func ProvideHealthCheckController(db *sqlx.DB, repo db.ModeratorUserRepository, collector service.ScraperStatCollector) []web.Route {
	controller := &HealthCheckController{
		db:        db,
		repo:      repo,
		collector: collector,
	}

	return controller.GetRoutes()
}

func (handler *HealthCheckController) GetRoutes() []web.Route {
	return []web.Route{
		web.NewRequestsHandler[HealthCheckController](handler.GetHealthCheckReport, "/health", []string{"GET"}),
	}
}

// Get the current health check report for the web service
//
// Path: `GET /api/v1/health`
//
// @Summary			Get the current health check report
// @Description		Endpoint for querying the most recent health check report for the web service
// @Tags			Health
// @Produce			json
// @Success			200 {object} dto.HealthCheckDTO
// @Router			/api/v1/health [get]
func (handler *HealthCheckController) GetHealthCheckReport(details *web.HttpRequestDetails[web.EmptyBody]) (int, interface{}, error) {
	hasAdmin, _ := handler.repo.HasAdministrator().Query()
	report := dto.HealthCheckDTO{
		DatabaseUp:       handler.db.Ping() == nil,
		IsInitialized:    hasAdmin != nil && *hasAdmin,
		LastScrapeResult: handler.collector.GetResults(),
	}

	return http.StatusOK, report, nil
}
