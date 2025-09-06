package bg

import (
	"encoding/json"
	"fmt"
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

const BENU_ENDPOINT = "https://www.benu.ee/leia-apteek"

type BenuScraper struct {
	repo       db.PharmacyRepository
	httpClient utils.HttpClient
	logger     zerolog.Logger
}

func ProvideBenuScraper(repo db.PharmacyRepository, client utils.HttpClient) Scraper {
	return BenuScraper{
		repo:       repo,
		httpClient: client,
		logger:     utils.GetLogger("BG"),
	}
}

type benuPharmacy struct {
	ID        int64  `json:"ID"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Region    string `json:"region"`
	Address   string `json:"address"`
	PostCode  string `json:"postCode"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	ModTime   string `json:"modTime"`
}

func (src *benuPharmacy) mapToPharmacy(dst *entity.Pharmacy, newTS time.Time, logger *zerolog.Logger) error {
	dst.PharmacyID = src.ID
	dst.Chain = string(entity.CHAIN_BENU)
	dst.County = src.Region
	dst.PostalCode = src.PostCode
	dst.Email = src.Email
	dst.PhoneNumber = fmt.Sprintf("+372%s", src.Phone)
	dst.ModTime = types.Time(newTS)
	lat, err := strconv.ParseFloat(src.Latitude, 32)
	if err != nil {
		logger.Error().Msgf("Failed to extract BENU pharmacy latitude: invalid value %s", src.Latitude)
		return fmt.Errorf("failed to extract pharmacy latitude: invalid value %s", src.Latitude)
	}

	lng, err := strconv.ParseFloat(src.Longitude, 32)
	if err != nil {
		logger.Error().Msgf("Failed to extract BENU pharmacy longitude: invalid value %s", src.Longitude)
		return fmt.Errorf("failed to extract pharmacy longitude: invalid value %s", src.Longitude)
	}
	dst.Latitude = float32(lat)
	dst.Longitude = float32(lng)

	// extracting address information with regex
	re := regexp.MustCompile(`^(.*?) *- *(.*?)( *- *(.*?))?( *- *(.*))?$`)
	groups := re.FindStringSubmatch(src.Address)

	// 7 groups means that the address also contains a district
	// 5 groups means no district but it has a city
	// 3 groups means no district and no city (can be extracted from the address part)
	if len(groups) == 7 {
		if groups[6] != "" {
			dst.City = groups[1]
			dst.Name = groups[4]
			dst.Address = fmt.Sprintf("%s, %s", groups[6], groups[2])
		} else if groups[4] != "" {
			dst.City = groups[1]
			dst.Name = groups[2]
			dst.Address = groups[4]
		} else if groups[2] != "" {
			dst.Name = groups[1]
			data := strings.Split(groups[2], ",")
			if len(data) == 2 {
				dst.Address = strings.Trim(data[0], " ")
				dst.City = strings.Trim(data[1], " ")
			}
		}
	}

	return nil
}

func (scraper *BenuScraper) createEntitiesFromJson(data string) ([]entity.Pharmacy, error) {
	var pharmacies map[string]benuPharmacy
	err := json.Unmarshal([]byte(data), &pharmacies)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to unmarshal BENU pharmacy json: %v", err)
		return nil, fmt.Errorf("failed to unmarshal BENU pharmacy json")
	}

	existing, err := scraper.repo.FindPharmaciesByChain(entity.CHAIN_BENU).QueryAll()
	if err != nil {
		scraper.logger.Error().Msgf("Failed to query existing BENU pharmacies in the database: %v", err)
		return nil, fmt.Errorf("failed to query existing BENU pharmacies in the database")
	}

	ret := make([]entity.Pharmacy, 0)
	for _, pharmacy := range pharmacies {
		newTS, err := time.Parse("2006-01-02 15:04:05", pharmacy.ModTime)
		if err != nil {
			// maybe the timestamp is missing
			newTS = time.Now().UTC()
		}

		var existingPharmacy *entity.Pharmacy
		for i := range existing {
			if pharmacy.ID == existing[i].PharmacyID {
				existingPharmacy = &existing[i]
				break
			}
		}

		// if existing pharmacy was found, then we check if it should be updated based on the timestamps
		if existingPharmacy != nil && (time.Time(existingPharmacy.ModTime).UTC().UnixMilli() < newTS.UTC().UnixMilli()) {
			err := pharmacy.mapToPharmacy(existingPharmacy, newTS, &scraper.logger)
			if err != nil {
				continue
			}
			ret = append(ret, *existingPharmacy)
		} else if existingPharmacy == nil {
			var newPharmacy entity.Pharmacy
			err := pharmacy.mapToPharmacy(&newPharmacy, newTS, &scraper.logger)
			if err != nil {
				continue
			}
			ret = append(ret, newPharmacy)
		}
	}

	return ret, nil
}

func (scraper BenuScraper) Scrape() {
	req, err := http.NewRequest("GET", BENU_ENDPOINT, nil)
	if err != nil {
		scraper.logger.Error().Msg("Failed to create a new request for BENU scraper")
		return
	}
	req.Header.Set("User-Agent", USER_AGENT)
	resp, err := scraper.httpClient.Do(req)
	if err != nil {
		scraper.logger.Error().Msgf("Failed to make a request to %s: %v", BENU_ENDPOINT, err)
		return
	}

	// make sure that the server responded with status code 200
	if resp.StatusCode != 200 {
		scraper.logger.Error().Msgf("Benu endpoint responded with non-200 status code %d", resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		scraper.logger.Error().Msg("Failed to read response body from BENU endpoint request")
		return
	}

	script := soup.HTMLParse(string(body)).
		Find("main").
		Find("div", "class", "bnContainer").
		Find("script")

	if script.Error != nil {
		scraper.logger.Error().Msg("Failed to extract script tag from BENU website's HTML body")
		return
	}

	txt := script.Text()
	re := regexp.MustCompile(`(?m)^.*?pharmacies = ({.+}).*$`)
	data := re.FindStringSubmatch(txt)[1]
	pharmacies, err := scraper.createEntitiesFromJson(data)

	if err != nil {
		scraper.logger.Error().Msgf("Failed to read pharmacy data from json: %v", err)
		return
	}

	scraper.repo.StoreAll(pharmacies)
}
