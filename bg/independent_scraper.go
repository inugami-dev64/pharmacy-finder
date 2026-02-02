package bg

import (
	"encoding/json"
	"hash/crc64"
	"io"
	"maps"
	"pharmafinder"
	"pharmafinder/db"
	"pharmafinder/db/entity"
	"pharmafinder/service"
	"pharmafinder/utils"

	"github.com/rs/zerolog"
)

// "Scrapes" independent pharmacies from embedded json file
type IndependentScraper struct {
	repo      db.PharmacyRepository
	collector service.ScraperStatCollector
	logger    zerolog.Logger
}

func ProvideIndependentScraper(repo db.PharmacyRepository, collector service.ScraperStatCollector) Scraper {
	return &IndependentScraper{
		repo:      repo,
		collector: collector,
		logger:    utils.GetLogger("BG"),
	}
}

func (scraper *IndependentScraper) Scrape() {
	success := false
	chains := make(map[string]bool)
	defer func() {
		for v := range maps.Keys(chains) {
			scraper.collector.CollectScrapeResult(entity.PharmacyChain(v), success)
		}
	}()

	// Load the embedded independent pharmacies json
	f, err := pharmafinder.PharmacyJSON.Open("db/independent-pharmacies.json")
	if err != nil {
		scraper.logger.Error().Msgf("Failed to open embedded db/independent-pharmacies.json file: %v", err)
		return
	}

	jsonBytes, err := io.ReadAll(f)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to read data from embedded file: %v", err)
		return
	}

	var pharmacies []entity.Pharmacy
	err = json.Unmarshal(jsonBytes, &pharmacies)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to unmarshal independent pharmacy json")
		return
	}

	var toStore []entity.Pharmacy
	for _, pharmacy := range pharmacies {
		chains[pharmacy.Chain] = true
		pharmacy.PharmacyID = int64(crc64.Checksum([]byte(pharmacy.Name), crc64Table))
		resps, err := scraper.repo.FindPharmacyByChainAndPharmacyID(pharmacy.PharmacyID, entity.PharmacyChain(pharmacy.Chain)).QueryAll()
		if err != nil {
			scraper.logger.Error().Msgf("Failed to query for existing pharmacies: %v", err)
			return
		}
		if len(resps) == 0 {
			toStore = append(toStore, pharmacy)
		}
	}

	if len(toStore) > 0 {
		err = scraper.repo.StoreAll(toStore)
		if err != nil {
			scraper.logger.Error().Msgf("Failed to persist independent pharmacies to the database: %v", err)
		}
	}

	success = true
}
