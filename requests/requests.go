package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"yet-another-covid-map-api/casecount"
)

func parseURLQuery(URL *url.URL, key string) string {
	keys, ok := URL.Query()[key]
	var result string
	if ok && len(keys) > 0 {
		result = keys[0]
	}
	return result
}

func checkForCaseCountError(err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}
	http.Error(w, err.Error(), http.StatusBadRequest)
	return true
}

// GetCaseCounts : logic when /cases endpoint is called. Returns all aggregated confirmed cases/death counts between from and to dates in the query
func GetCaseCounts(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())
	from := parseURLQuery(r.URL, "from")
	to := parseURLQuery(r.URL, "to")
	country := parseURLQuery(r.URL, "country")
	var response []byte
	var err error
	if parseURLQuery(r.URL, "aggregateCountries") == "true" {
		caseCounts, caseCountsErr := casecount.GetCountryCaseCounts(from, to, country)
		if checkForCaseCountError(caseCountsErr, w) {
			return
		}
		response, err = json.Marshal(caseCounts)
	} else {
		caseCounts, caseCountsErr := casecount.GetCaseCounts(from, to, country)
		if checkForCaseCountError(caseCountsErr, w) {
			return
		}
		response, err = json.Marshal(caseCounts)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
