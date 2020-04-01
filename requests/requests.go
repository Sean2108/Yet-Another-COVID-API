package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"yet-another-covid-map-api/casecount"
)

func parseURL(URL *url.URL, getAbbreviation bool) (string, string, string, bool) {
	from := parseURLQuery(URL, "from")
	to := parseURLQuery(URL, "to")
	country := parseURLQuery(URL, "country")
	var countryLookupFuncToCall func(string) (string, bool)
	if getAbbreviation {
		countryLookupFuncToCall = getAbbreviationFromCountry
	} else {
		countryLookupFuncToCall = getCountryFromAbbreviation
	}
	if countryFromAbbr, ok := countryLookupFuncToCall(country); ok {
		country = countryFromAbbr
	}
	aggregateCountries := parseURLQuery(URL, "aggregateCountries") == "true"
	return from, to, country, aggregateCountries
}

func parseURLQuery(URL *url.URL, key string) string {
	keys, ok := URL.Query()[key]
	var result string
	if ok && len(keys) > 0 {
		result = keys[0]
	}
	return result
}

func getCaseCountsResponse(from string, to string, country string, aggregateCountries bool) ([]byte, error, error) {
	if aggregateCountries {
		caseCounts, caseCountsErr := casecount.GetCountryCaseCounts(from, to, country)
		if caseCountsErr != nil {
			return nil, nil, caseCountsErr
		}
		response, err := json.Marshal(caseCounts)
		return response, err, nil
	}
	caseCounts, caseCountsErr := casecount.GetCaseCounts(from, to, country)
	if caseCountsErr != nil {
		return nil, nil, caseCountsErr
	}
	response, err := json.Marshal(caseCounts)
	return response, err, nil
}

// GetCaseCounts : logic when /cases endpoint is called. Returns all aggregated confirmed cases/death counts between from and to dates in the query
func GetCaseCounts(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())
	from, to, country, aggregateCountries := parseURL(r.URL, false)
	response, err, caseCountsErr := getCaseCountsResponse(from, to, country, aggregateCountries)
	if caseCountsErr != nil {
		http.Error(w, caseCountsErr.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
