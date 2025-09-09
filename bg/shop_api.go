package bg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pharmafinder/db/entity"
	"pharmafinder/types"
	"pharmafinder/utils"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// I noticed that Apotheka and Südameapteek
// use suspiciously similar APIs for their sites
// thus I decided to write unified scraper logic
// here to avoid code duplication

type shop struct {
	ShopID            string  `json:"shop_id"`
	Name              string  `json:"name"`
	City              string  `json:"city"`
	County            string  `json:"districtName"`
	Address           string  `json:"address"`
	Email             string  `json:"email"`
	Phone             string  `json:"phone"`
	UpdatedAt         string  `json:"updated_at"`
	LocationLatitude  float32 `json:"location_latitude"`
	LocationLongitude float32 `json:"location_longitude"`
}

type shops struct {
	TotalRecords int    `json:"totalRecords"`
	Items        []shop `json:"items"`
}

// Fetch Apotheka or Südameapteek pharmacies
// from provided API endpoint
func fetchShops(url string, client utils.HttpClient, logger *zerolog.Logger) (*shops, error) {
	logger.Debug().Msgf("Fetching shops from API endpoint %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error().Msgf("Failed to create a GET request object to %s: %v", url, err)
		return nil, fmt.Errorf("failed to request shops from %s", url)
	}

	req.Header.Set("User-Agent", USER_AGENT)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Msgf("Failed to make a request to %s: %v", url, err)
		return nil, fmt.Errorf("failed to make an HTTP request to %s: %v", url, err)
	}

	// make sure that the server responded with status code 200
	if resp.StatusCode != 200 {
		logger.Error().Msgf("API responded with non-200 status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("shop API responded with unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Msg("Failed to read response body from shop API")
		return nil, fmt.Errorf("failed to read response body from shop API")
	}

	var pharmacies shops
	err = json.Unmarshal(body, &pharmacies)
	if err != nil {
		logger.Error().Msgf("Failed to unmarshal shop API response body: %v", err)
		return nil, fmt.Errorf("failed to unmarshal shop API response: %v", err)
	}

	return &pharmacies, nil
}

func mapShopsToPharmacies(pharmacyShops *shops, chain entity.PharmacyChain, logger *zerolog.Logger, client utils.HttpClient) []entity.Pharmacy {
	pharmacies := make([]entity.Pharmacy, 0)
	for i := range pharmacyShops.Items {
		pharmacyID, err := strconv.ParseInt(pharmacyShops.Items[i].ShopID, 10, 64)
		if err != nil {
			logger.Warn().Msgf("Failed to extract pharmacy ID for %s pharmacy %s, skipping", chain, pharmacyShops.Items[i].Name)
			continue
		}

		var pharmacy entity.Pharmacy
		pharmacy.PharmacyID = pharmacyID
		pharmacy.Chain = string(chain)
		pharmacy.Name = pharmacyShops.Items[i].Name

		addressParts := strings.Split(pharmacyShops.Items[i].Address, ",")
		if len(addressParts) > 0 {
			pharmacy.Address = strings.Trim(addressParts[0], " ")
		} else {
			pharmacy.Address = pharmacyShops.Items[i].Address
		}

		pharmacy.City = pharmacyShops.Items[i].City
		pharmacy.County = pharmacyShops.Items[i].County

		pharmacy.PostalCode = fetchOmnivaZipCode(pharmacyShops.Items[i].Address, client, logger)
		pharmacy.Email = pharmacyShops.Items[i].Email
		pharmacy.PhoneNumber = fmt.Sprintf("+372%s", pharmacyShops.Items[i].Phone)
		ts, err := time.Parse("2006-01-02 15:04:05", pharmacyShops.Items[i].UpdatedAt)
		if err != nil {
			logger.Warn().Msgf("Failed to extract modification timestamp for pharmacy '%s'", pharmacy.Name)
		} else {
			pharmacy.ModTime = types.Time(ts)
		}

		pharmacy.Latitude = pharmacyShops.Items[i].LocationLatitude
		pharmacy.Longitude = pharmacyShops.Items[i].LocationLongitude

		pharmacies = append(pharmacies, pharmacy)
	}

	return pharmacies
}
