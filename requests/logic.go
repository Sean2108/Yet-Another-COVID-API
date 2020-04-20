package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"yet-another-covid-map-api/casecount"
	"yet-another-covid-map-api/dateformat"
	"yet-another-covid-map-api/news"
)

type writer interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

func parseURL(URL *url.URL, getAbbreviation bool, dateFormat string) (string, string, string, bool, bool, bool, bool) {
	from := parseURLQuery(URL, "from")
	to := parseURLQuery(URL, "to")
	country := parseURLQuery(URL, "country")

	from, ok := dateformat.FormatDate(dateFormat, from)
	if !ok {
		return "", "", "", false, false, false, false
	}
	to, ok = dateformat.FormatDate(dateFormat, to)
	if !ok {
		return "", "", "", false, false, false, false
	}

	var countryLookupFuncToCall func(string) (string, bool)
	if getAbbreviation {
		countryLookupFuncToCall = getAbbreviationFromCountry
	} else {
		countryLookupFuncToCall = getCountryFromAbbreviation
	}
	if countryFromAbbr, ok := countryLookupFuncToCall(country); ok {
		country = countryFromAbbr
	}
	aggregateCountries := parseURLQuery(URL, "aggregatecountries") == "true"
	perDay := parseURLQuery(URL, "perday") == "true"
	worldTotal := parseURLQuery(URL, "worldtotal") == "true"

	return from, to, country, aggregateCountries, perDay, worldTotal, true
}

func parseURLQuery(URL *url.URL, key string) string {
	query := URL.Query()
	for k, v := range query {
		if strings.ToLower(k) == key && len(v) > 0 {
			return v[0]
		}
	}
	return ""
}

func getCaseCountsResponse(from string, to string, country string, aggregateCountries bool, perDay bool, worldTotal bool) ([]byte, error, error) {
	if worldTotal {
		caseCounts, caseCountsErr := casecount.GetWorldCaseCounts(from, to)
		response, err := json.Marshal(caseCounts)
		return response, err, caseCountsErr
	}
	if perDay {
		if aggregateCountries {
			caseCounts, caseCountsErr := casecount.GetCountryCaseCountsWithDayData(from, to, country)
			response, err := json.Marshal(caseCounts)
			return response, err, caseCountsErr
		}
		caseCounts, caseCountsErr := casecount.GetCaseCountsWithDayData(from, to, country)
		response, err := json.Marshal(caseCounts)
		return response, err, caseCountsErr
	}
	if aggregateCountries {
		caseCounts, caseCountsErr := casecount.GetCountryCaseCounts(from, to, country)
		response, err := json.Marshal(caseCounts)
		return response, err, caseCountsErr
	}
	caseCounts, caseCountsErr := casecount.GetCaseCounts(from, to, country)
	response, err := json.Marshal(caseCounts)
	return response, err, caseCountsErr
}

func getNewsForCountryResponse(from string, to string, country string, _ bool, _ bool, _ bool) ([]byte, error, error) {
	articles, newsErr := news.GetNews(from, to, country)
	response, err := json.Marshal(articles)
	return response, err, newsErr
}

func getResponse(getDataFn func(from string, to string, country string, aggregateCountries bool, perDay bool, worldTotal bool) ([]byte, error, error), w writer, URL *url.URL, getCountryAbbreviation bool) {
	log.Println(URL.String())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	from, to, country, aggregateCountries, perDay, worldTotal, ok := parseURL(URL, getCountryAbbreviation, dateformat.CasesDateFormat)
	if !ok {
		http.Error(w, "Date format is not recognised, please use either YYYY-MM-DD, YYYY/MM/DD, MM-DD-YY or MM/DD/YY", http.StatusBadRequest)
		return
	}
	response, jsonErr, internalErr := getDataFn(from, to, country, aggregateCountries, perDay, worldTotal)
	if internalErr != nil {
		http.Error(w, internalErr.Error(), http.StatusBadRequest)
		return
	}
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
