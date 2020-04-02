package utils

import "net/http"

// HTTPClient : Interface to mock net/http client
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}
