package dto

import (
	"pharmafinder/db/entity"
	"pharmafinder/types"
)

type PharmacyScraperResultDTO struct {
	Chain                         entity.PharmacyChain `json:"chain"`
	LastSuccessfulScrapeTimestamp *types.Time          `json:"lastSuccessfulScrapeTimestamp,omitempty"`
	Timestamp                     types.Time           `json:"timestamp"`
	Success                       bool                 `json:"success"`
}

type HealthCheckDTO struct {
	DatabaseUp       bool                       `json:"dbUp"`
	IsInitialized    bool                       `json:"initialized"`
	LastScrapeResult []PharmacyScraperResultDTO `json:"lastScrapeResults"`
}
