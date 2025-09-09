package bg

import (
	"pharmafinder/db"
	"pharmafinder/db/entity"
	"pharmafinder/utils"

	"github.com/rs/zerolog"
)

const APOTHEKA_ENDPOINT = "https://www.apotheka.ee/shops/shop/shops"

type ApothekaScraper struct {
	repo       db.PharmacyRepository
	httpClient utils.HttpClient
	logger     zerolog.Logger
}

func ProvideApothekaScraper(repo db.PharmacyRepository, client utils.HttpClient) Scraper {
	return &ApothekaScraper{
		repo:       repo,
		httpClient: client,
		logger:     utils.GetLogger("BG"),
	}
}

func (scraper *ApothekaScraper) Scrape() {
	pharmacies, err := fetchShops(APOTHEKA_ENDPOINT, scraper.httpClient, &scraper.logger)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to fetch Apotheka pharmacies: %v", err)
		return
	}

	apothekaPharmacies := mapShopsToPharmacies(pharmacies, entity.CHAIN_APOTHEKA, &scraper.logger, scraper.httpClient)
	err = scraper.repo.StoreAll(apothekaPharmacies)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to persist Apotheka pharmacies: %v", err)
	}
}
