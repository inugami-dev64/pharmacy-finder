package bg

import (
	"pharmafinder/db"
	"pharmafinder/db/entity"
	"pharmafinder/service"
	"pharmafinder/utils"
	"time"

	"github.com/rs/zerolog"
)

const APOTHEKA_ENDPOINT = "https://www.apotheka.ee/shops/shop/shops"

type ApothekaScraper struct {
	repo       db.PharmacyRepository
	httpClient utils.HttpClient
	collector  service.ScraperStatCollector
	logger     zerolog.Logger
}

func ProvideApothekaScraper(repo db.PharmacyRepository, client utils.HttpClient, collector service.ScraperStatCollector) Scraper {
	return &ApothekaScraper{
		repo:       repo,
		httpClient: client,
		collector:  collector,
		logger:     utils.GetLogger("BG"),
	}
}

func (scraper *ApothekaScraper) Scrape() {
	success := false
	defer func() { scraper.collector.CollectScrapeResult(entity.CHAIN_APOTHEKA, success) }()

	scraper.logger.Info().Msg("Scraping Apotheka pharmacy locations...")
	existingPharmacies, err := scraper.repo.FindPharmaciesByChain(entity.CHAIN_APOTHEKA).QueryAll()
	if err != nil {
		scraper.logger.Error().Msgf("Failed to query existing Apotheka pharmacies: %v", err)
		return
	}

	pharmacies, err := fetchShops(APOTHEKA_ENDPOINT, scraper.httpClient, &scraper.logger)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to fetch Apotheka pharmacies: %v", err)
		return
	}

	pharmaciesToSave := make([]entity.Pharmacy, 0)
	apothekaPharmacies := mapShopsToPharmacies(pharmacies, entity.CHAIN_APOTHEKA, &scraper.logger, scraper.httpClient)
	for i := range apothekaPharmacies {
		var existingPharmacy *entity.Pharmacy
		for _, existing := range existingPharmacies {
			if apothekaPharmacies[i].PharmacyID == existing.PharmacyID {
				existingPharmacy = &existing
				break
			}
		}

		if existingPharmacy != nil && time.Time(existingPharmacy.ModTime).UTC().UnixMilli() < time.Time(apothekaPharmacies[i].ModTime).UTC().UnixMilli() {
			apothekaPharmacies[i].ID = existingPharmacy.ID
			pharmaciesToSave = append(pharmaciesToSave, apothekaPharmacies[i])
		} else if existingPharmacy == nil {
			pharmaciesToSave = append(pharmaciesToSave, apothekaPharmacies[i])
		}
	}

	err = scraper.repo.StoreAll(pharmaciesToSave)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to persist Apotheka pharmacies: %v", err)
		return
	}

	success = true
}
