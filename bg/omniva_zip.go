package bg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pharmafinder/utils"
	"strings"

	"github.com/rs/zerolog"
)

const OMNIVA_ZIP_CODE_ENDPOINT = "https://www.omniva.ee/wp-json/custom/v1/omniva-zip-search?search=%s"

type zipcodeResponse struct {
	Addresses []struct {
		Address string `json:"address"`
		ZipCode string `json:"zipCode"`
	} `json:"addresses"`
}

func fetchOmnivaZipCode(address string, client utils.HttpClient, logger *zerolog.Logger) string {
	escapedAddress := strings.ReplaceAll(url.QueryEscape(address), "%20", "+")
	url := fmt.Sprintf(OMNIVA_ZIP_CODE_ENDPOINT, escapedAddress)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error().Msg("Failed to create a new request object for Omniva zip code API")
		return ""
	}

	req.Header.Set("User-Agent", USER_AGENT)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Msgf("Failed to make a request to %s: %v", url, err)
		return ""
	}

	// make sure that the server responded with status code 200
	if resp.StatusCode != 200 {
		logger.Error().Msgf("Omniva zipcode API responded with non-200 status code %d", resp.StatusCode)
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Msg("Failed to read response body from Omniva zipcode API response")
		return ""
	}

	var zipcode zipcodeResponse
	err = json.Unmarshal(body, &zipcode)
	if err != nil {
		logger.Error().Msgf("Failed to unmarshal Omniva zipcode API response: %v", err)
		return ""
	}

	if len(zipcode.Addresses) > 0 {
		return zipcode.Addresses[0].ZipCode
	}

	return ""
}
