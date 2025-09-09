package bg

import (
	"pharmafinder/db"
	"pharmafinder/db/entity"
	"pharmafinder/utils"

	"github.com/rs/zerolog"
)

const SYDAMEAPTEEK_ENDPOINT = "https://www.sudameapteek.ee/shops/shop/shops"

type SydameapteekScraper struct {
	repo       db.PharmacyRepository
	httpClient utils.HttpClient
	logger     zerolog.Logger
}

func ProvideSydameapteekScraper(repo db.PharmacyRepository, client utils.HttpClient) Scraper {
	return &SydameapteekScraper{
		repo:       repo,
		httpClient: client,
		logger:     utils.GetLogger("BG"),
	}
}

func (scraper *SydameapteekScraper) Scrape() {
	pharmacies, err := fetchShops(APOTHEKA_ENDPOINT, scraper.httpClient, &scraper.logger)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to fetch Südameapteek pharmacies: %v", err)
		return
	}

	apothekaPharmacies := mapShopsToPharmacies(pharmacies, entity.CHAIN_APOTHEKA, &scraper.logger, scraper.httpClient)
	err = scraper.repo.StoreAll(apothekaPharmacies)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to persist Südameapteek pharmacies: %v", err)
	}
}
