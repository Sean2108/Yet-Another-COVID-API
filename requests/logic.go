package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"yet-another-covid-map-api/casecount"
	"yet-another-covid-map-api/dateformat"
	"yet-another-covid-map-api/news"
	"yet-another-covid-map-api/utils"
)

type writer interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

func parseURL(URL *url.URL, dateFormat string) (string, string, string, bool, bool, bool, error) {
	from := parseURLQuery(URL, "from")
	to := parseURLQuery(URL, "to")
	country := parseURLQuery(URL, "country")

	from, fromOk := dateformat.FormatDate(dateFormat, from)
	to, toOk := dateformat.FormatDate(dateFormat, to)
	if !fromOk || !toOk {
		return "", "", "", false, false, false, errors.New("Date format is not recognised, please use either YYYY-MM-DD, YYYY/MM/DD, MM-DD-YY or MM/DD/YY")
	}

	if country != "" {
		if countryFromAbbr, ok := utils.GetAbbreviationFromCountry(country); ok {
			country = countryFromAbbr
		} else {
			return "", "", "", false, false, false, fmt.Errorf("Country %s not found, did you mean: %s?", country, countryFromAbbr)
		}
	}
	aggregateCountries := isStringTrue(parseURLQuery(URL, "aggregatecountries"))
	perDay := isStringTrue(parseURLQuery(URL, "perday"))
	worldTotal := isStringTrue(parseURLQuery(URL, "worldtotal"))

	return from, to, country, aggregateCountries, perDay, worldTotal, nil
}

func isStringTrue(str string) bool {
	if val, err := strconv.ParseBool(str); err == nil {
		return val
	}
	return false
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
	from, to, country, aggregateCountries, perDay, worldTotal, err := parseURL(URL, dateformat.CasesDateFormat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
