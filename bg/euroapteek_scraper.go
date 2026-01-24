package bg

import (
	"encoding/json"
	"fmt"
	"hash/crc64"
	"io"
	"net/http"
	"pharmafinder/db"
	"pharmafinder/db/entity"
	"pharmafinder/types"
	"pharmafinder/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/rs/zerolog"
)

const EUROAPTEEK_WEBSITE = "https://www.euroapteek.ee/apteegid"

type EuroapteekScraper struct {
	repo       db.PharmacyRepository
	httpClient utils.HttpClient
	logger     zerolog.Logger
}

type euroapteekPharmacy struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	City        string `json:"city"`
	County      string `json:"country"`
	Latitude    string `json:"lat"`
	Longitude   string `json:"lng"`
}

var crc64Table *crc64.Table = crc64.MakeTable(crc64.ISO)

func ProvideEuroapteekScraper(repo db.PharmacyRepository, client utils.HttpClient) Scraper {
	return &EuroapteekScraper{
		repo:       repo,
		httpClient: client,
		logger:     utils.GetLogger("BG"),
	}
}

func (scraper EuroapteekScraper) mapToPharmacies(existingPharmacies []entity.Pharmacy, scrapedPharmacies []euroapteekPharmacy) []entity.Pharmacy {
	pharmacies := make([]entity.Pharmacy, 0)
	for _, scraped := range scrapedPharmacies {
		var existingPharmacy *entity.Pharmacy
		pharmacyID := crc64.Checksum([]byte(scraped.Name), crc64Table)
		for i := range existingPharmacies {
			if existingPharmacies[i].PharmacyID == int64(pharmacyID) {
				existingPharmacy = &existingPharmacies[i]
				break
			}
		}

		if existingPharmacy == nil {
			var pharmacy entity.Pharmacy
			pharmacy.PharmacyID = int64(crc64.Checksum([]byte(scraped.Name), crc64Table))
			pharmacy.Chain = string(entity.CHAIN_EUROAPTEEK)
			pharmacy.Name = scraped.Name
			pharmacy.Address = scraped.Address
			pharmacy.City = scraped.City
			pharmacy.County = scraped.County
			// Commented out for testing reasons (this is extremely slow query)
			pharmacy.PostalCode = fetchOmnivaZipCode(fmt.Sprintf("%s, %s, %s", pharmacy.Address, pharmacy.City, pharmacy.County), scraper.httpClient, &scraper.logger)
			pharmacy.ModTime = types.Time(time.UnixMilli(0))

			// extract coordinates (lat, lng)
			lat, err := strconv.ParseFloat(scraped.Latitude, 32)
			if err != nil {
				scraper.logger.Error().Msgf("Failed to extract latitude for Euroapteek pharmacy %s: %v", pharmacy.Name, err)
				continue
			}
			pharmacy.Latitude = float32(lat)

			lng, err := strconv.ParseFloat(scraped.Longitude, 32)
			if err != nil {
				scraper.logger.Error().Msgf("Failed to extract longitude for Euroapteek pharmacy %s: %v", pharmacy.Name, err)
				continue
			}
			pharmacy.Longitude = float32(lng)

			// extract and sanitize phone number
			re := regexp.MustCompile(`(\+372)? *([\d ]+)`)
			groups := re.FindStringSubmatch(scraped.PhoneNumber)
			if len(groups) == 3 {
				pharmacy.PhoneNumber = fmt.Sprintf("+372%s", strings.ReplaceAll(groups[2], " ", ""))
			} else {
				scraper.logger.Warn().Msgf("Failed to extract phone number for Euroapteek pharmacy %s", pharmacy.Name)
			}

			pharmacies = append(pharmacies, pharmacy)
		}
	}

	return pharmacies
}

func (scraper *EuroapteekScraper) Scrape() {
	scraper.logger.Info().Msg("Scraping Euroapteek pharmacy locations...")

	existingPharmacies, err := scraper.repo.FindPharmaciesByChain(entity.CHAIN_EUROAPTEEK).QueryAll()
	if err != nil {
		scraper.logger.Error().Msgf("Failed to query existing Euroapteek pharmacies: %v", err)
		return
	}

	req, err := http.NewRequest("GET", EUROAPTEEK_WEBSITE, nil)
	if err != nil {
		scraper.logger.Error().Msg("Failed to create a request object for Euroapteek API")
		return
	}

	resp, err := scraper.httpClient.Do(req)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to make a request to Euroapteek API: %v", err)
		return
	}

	if resp.StatusCode != 200 {
		scraper.logger.Error().Msgf("Euroapteek API responded with non-200 status code %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to read response body from Euroapteek API: %v", err)
		return
	}

	var apiPharmacies []euroapteekPharmacy
	err = json.Unmarshal(body, &apiPharmacies)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to unmarshal JSON from Euroapteek API response: %v", err)
		return
	}

	pharmacies := scraper.mapToPharmacies(existingPharmacies, apiPharmacies)
	err = scraper.repo.StoreAll(pharmacies)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to persist Euroapteek pharmacies: %v", err)
	}
}
