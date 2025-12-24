package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"pharmafinder/utils"
	"slices"
	"strings"

	"github.com/rs/zerolog"
)

const RECAPTCHA_VERIFY_ENDPOINT = "https://www.google.com/recaptcha/api/siteverify"

type RecaptchaVerifier interface {
	// Verify if provided grecaptcha response is valid by making an
	// API request to Google
	//
	// If the challenge has been successfully completed by the user, returns true
	// false otherwise
	Verify(response string) bool
}

type recaptchaVerificationRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
}

type recaptchaVerificationResponse struct {
	Success     bool   `json:"success"`
	ChallengeTS string `json:"challenge_ts"`
	Hostname    string `json:"hostname"`
}

type RecaptchaVerifierImpl struct {
	client utils.HttpClient
	logger zerolog.Logger
}

func ProvideRecaptchaVerifier(client utils.HttpClient) RecaptchaVerifierImpl {
	return RecaptchaVerifierImpl{
		client: client,
		logger: utils.GetLogger("SERVICE"),
	}
}

func (verifier RecaptchaVerifierImpl) Verify(response string) bool {
	reqBody := recaptchaVerificationRequest{
		Secret:   utils.Getenv("RECAPTCHA_SECRET", ""),
		Response: response,
	}

	reqJson, err := json.Marshal(reqBody)
	if err != nil {
		verifier.logger.Error().Msgf("Failed to marshal request body for reCaptcha verification: %v", err)
		return false
	}

	req, err := http.NewRequest("POST", RECAPTCHA_VERIFY_ENDPOINT, bytes.NewReader(reqJson))
	if err != nil {
		verifier.logger.Error().Msgf("Failed to create a new http.Request instance for verifying reCaptcha response: %v", err)
		return false
	}

	resp, err := verifier.client.Do(req)
	if err != nil {
		verifier.logger.Error().Msgf("Failed to make a request to %s: %v", RECAPTCHA_VERIFY_ENDPOINT, err)
		return false
	}

	if resp.StatusCode != 200 {
		verifier.logger.Error().Msgf("ReCaptcha verification endpoint returned non-200 status code %d", resp.StatusCode)
		return false
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		verifier.logger.Error().Msgf("Failed to read reCaptcha verification endpoint response: %v", err)
		return false
	}

	var grResp recaptchaVerificationResponse
	json.Unmarshal(respBytes, &grResp)

	allowedDomains := strings.Split(utils.Getenv("ALLOWED_DOMAINS", ""), ",")
	return grResp.Success && slices.Contains(allowedDomains, grResp.Hostname)
}
