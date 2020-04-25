package utils

import (
	"encoding/csv"
	"net/http"
)

// HTTPClient : Interface to mock net/http client
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

func ReadCSVFromURL(client HTTPClient, url string) ([][]string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}
