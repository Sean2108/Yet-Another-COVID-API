package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"yet-another-covid-map-api/casecount"
	"yet-another-covid-map-api/dateformat"
	"yet-another-covid-map-api/news"
)

type writer interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

func parseURL(URL *url.URL, getAbbreviation bool, dateFormat string) (string, string, string, bool, bool, bool) {
	from := parseURLQuery(URL, "from")
	to := parseURLQuery(URL, "to")
	country := parseURLQuery(URL, "country")

	from, ok := dateformat.FormatDate(dateFormat, from)
	if !ok {
		return "", "", "", false, false, false
	}
	to, ok = dateformat.FormatDate(dateFormat, to)
	if !ok {
		return "", "", "", false, false, false
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
	aggregateCountries := parseURLQuery(URL, "aggregateCountries") == "true"
	perDay := parseURLQuery(URL, "perDay") == "true"

	return from, to, country, aggregateCountries, perDay, true
}

func parseURLQuery(URL *url.URL, key string) string {
	keys, ok := URL.Query()[key]
	var result string
	if ok && len(keys) > 0 {
		result = keys[0]
	}
	return result
}

func getCaseCountsResponse(from string, to string, country string, aggregateCountries bool, perDay bool) ([]byte, error, error) {
	if perDay {
		if aggregateCountries {
			CaseCounts, caseCountsErr := casecount.GetCountryCaseCountsWithDayData(from, to, country)
			response, err := json.Marshal(CaseCounts)
			return response, err, caseCountsErr
		}
		CaseCounts, caseCountsErr := casecount.GetCaseCountsWithDayData(from, to, country)
		response, err := json.Marshal(CaseCounts)
		return response, err, caseCountsErr
	}
	if aggregateCountries {
		CaseCounts, caseCountsErr := casecount.GetCountryCaseCounts(from, to, country)
		response, err := json.Marshal(CaseCounts)
		return response, err, caseCountsErr
	}
	CaseCounts, caseCountsErr := casecount.GetCaseCounts(from, to, country)
	response, err := json.Marshal(CaseCounts)
	return response, err, caseCountsErr
}

func getNewsForCountryResponse(from string, to string, country string, _ bool, _ bool) ([]byte, error, error) {
	articles, newsErr := news.GetNews(from, to, country)
	response, err := json.Marshal(articles)
	return response, err, newsErr
}

func getResponse(getDataFn func(from string, to string, country string, aggregateCountries bool, perDay bool) ([]byte, error, error), w writer, URL *url.URL) {
	log.Println(URL.String())
	from, to, country, aggregateCountries, perDay, ok := parseURL(URL, false, dateformat.CasesDateFormat)
	if !ok {
		http.Error(w, "Date format is not recognised, please use either YYYY-MM-DD, YYYY/MM/DD, MM-DD-YY or MM/DD/YY", http.StatusBadRequest)
		return
	}
	response, jsonErr, internalErr := getDataFn(from, to, country, aggregateCountries, perDay)
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
