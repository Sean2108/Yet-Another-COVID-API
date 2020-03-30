package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"yet-another-covid-map-api/casecount"
)

// CountryNews : TODO Will remove later when placeholder function getNewsForCountry is implemented
type CountryNews struct {
	News string
}

func parseURLQuery(URL *url.URL, key string) string {
	keys, ok := URL.Query()[key]
	var result string
	if ok && len(keys) > 0 {
		result = keys[0]
	}
	return result
}

// GetCaseCounts : logic when /cases endpoint is called. Returns all aggregated confirmed cases/death counts between from and to dates in the query
func GetCaseCounts(w http.ResponseWriter, r *http.Request) {
	from := parseURLQuery(r.URL, "from")
	to := parseURLQuery(r.URL, "to")
	if from == "" && to == "" {
		log.Println("GetCaseCounts query for all data")
	} else {
		log.Printf("GetCaseCounts query from: %s, to: %s\n", from, to)
	}
	caseCounts := casecount.GetCaseCounts(from, to)
	response, err := json.Marshal(caseCounts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// GetNewsForCountry : TODO, runs query to get all virus related news for a given country
func GetNewsForCountry(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["country"]
	if !ok || len(keys) < 1 {
		log.Printf("Url param 'country' is missing in request: %s", r.URL)
		http.Error(w, "Url param 'country' is missing in request", http.StatusBadRequest)
		return
	}
	countryNews := CountryNews{keys[0]}
	response, err := json.Marshal(countryNews)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
