package utils

import (
	"encoding/csv"
	"log"
	"net/http"
)

// HTTPClient : Interface to mock net/http client
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

func ReadCSVFromURL(client HTTPClient, url string) ([][]string, bool) {
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error was encountered getting the data from %s: %s\n", url, err.Error())
		return nil, false
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error was encountered reading the data from %s: %s\n", url, err.Error())
		return nil, false
	}

	return data, true
}
