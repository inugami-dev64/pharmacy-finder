package bg

import (
	"pharmafinder/db"
	"pharmafinder/db/entity"
	"pharmafinder/utils"
	"time"

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
	scraper.logger.Info().Msg("Scraping Südameapteek pharmacy locations...")
	existingPharmacies, err := scraper.repo.FindPharmaciesByChain(entity.CHAIN_SUDAMEAPTEEK).QueryAll()
	if err != nil {
		scraper.logger.Error().Msgf("Failed to query existing Südameapteek pharmacies: %v", err)
		return
	}

	pharmacies, err := fetchShops(SYDAMEAPTEEK_ENDPOINT, scraper.httpClient, &scraper.logger)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to fetch Südameapteek pharmacies: %v", err)
		return
	}

	pharmaciesToSave := make([]entity.Pharmacy, 0)
	sudameapteekPharmacies := mapShopsToPharmacies(pharmacies, entity.CHAIN_SUDAMEAPTEEK, &scraper.logger, scraper.httpClient)
	for i := range sudameapteekPharmacies {
		var existingPharmacy *entity.Pharmacy
		for _, existing := range existingPharmacies {
			if sudameapteekPharmacies[i].PharmacyID == existing.PharmacyID {
				existingPharmacy = &existing
				break
			}
		}

		if existingPharmacy != nil && time.Time(existingPharmacy.ModTime).UTC().UnixMilli() < time.Time(sudameapteekPharmacies[i].ModTime).UTC().UnixMilli() {
			sudameapteekPharmacies[i].ID = existingPharmacy.ID
			pharmaciesToSave = append(pharmaciesToSave, sudameapteekPharmacies[i])
		} else if existingPharmacy == nil {
			pharmaciesToSave = append(pharmaciesToSave, sudameapteekPharmacies[i])
		}
	}

	err = scraper.repo.StoreAll(pharmaciesToSave)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to persist Südameapteek pharmacies: %v", err)
	}
}
