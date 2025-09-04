package utils

import "net/http"

// Small abstraction interface for better DI reasons
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func ProvideHTTPClient() HttpClient {
	return http.DefaultClient
}
