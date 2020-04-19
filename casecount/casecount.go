package casecount

import (
	"log"
	"net/http"
	"time"

	"yet-another-covid-map-api/dateformat"
	"yet-another-covid-map-api/utils"
)

const (
	confirmedURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv"
	deathsURL    = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv"
)

var (
	// cache the query for getting all data for all states and all countries, because it is the most heavily used
	allAggregatedData          []CaseCountsAggregated
	allCountriesAggregatedData []CountryCaseCountsAggregated

	caseCountsCache        []CaseCounts
	worldCaseCountsCache   []CaseCount
	countryCaseCountsCache []CountryCaseCounts
	lastDate               time.Time
	firstDate              time.Time

	client utils.HTTPClient
)

func init() {
	client = &http.Client{}
}

// UpdateCaseCounts : Pull data from the John Hopkins CSV files on GitHub, store the result in a cache and also cache the aggregate data for the entire period
func UpdateCaseCounts() {
	log.Println("Updating case counts")
	caseCountsCache = nil
	confirmedData, deathsData := getData()
	headerRow := confirmedData[0]
	extractCaseCounts(headerRow, confirmedData, deathsData)
	setDateBoundariesAndAllAggregatedData(headerRow)
}

// GetCaseCountsWithDayData : get case counts for states but without aggregating the counts, so a list of days with number of confirmed cases and deaths on each day is returned
func GetCaseCountsWithDayData(from string, to string, country string) ([]CaseCounts, error) {
	if from == "" && to == "" && country == "" {
		log.Println("GetCaseCounts query for all data with per day information")
		return caseCountsCache, nil
	}
	log.Printf("GetCaseCountsWithDayData query from: %s, to: %s, country: %s\n", from, to, country)
	return filterCaseCounts(from, to, country)
}

// GetCountryCaseCountsWithDayData : get case counts for countries but without aggregating the counts, so a list of days with number of confirmed cases and deaths on each day is returned
func GetCountryCaseCountsWithDayData(from string, to string, country string) ([]CountryCaseCounts, error) {
	if from == "" && to == "" && country == "" {
		log.Println("GetCountryCaseCounts query for all data with per day information")
		return countryCaseCountsCache, nil
	}
	log.Printf("GetCountryCaseCountsWithDayData query from: %s, to: %s, country: %s\n", from, to, country)
	filtered, err := filterCaseCounts(from, to, country)
	return aggregateCountryDataFromCaseCounts(filtered), err
}

// GetCaseCounts : get case counts for all states between from date and to date. Return case counts for entire period if from and to dates are empty strings
func GetCaseCounts(from string, to string, country string) ([]CaseCountsAggregated, error) {
	if from == "" && to == "" && country == "" {
		log.Println("GetCaseCounts query for all data")
		return allAggregatedData, nil
	}
	log.Printf("GetCaseCounts query from: %s, to: %s, country: %s\n", from, to, country)
	return aggregateDataBetweenDates(from, to, country)
}

// GetCountryCaseCounts : get case counts for all countries between from date and to date. Return case counts for entire period if from and to dates are empty strings
func GetCountryCaseCounts(from string, to string, country string) ([]CountryCaseCountsAggregated, error) {
	if from == "" && to == "" && country == "" {
		log.Println("GetCountryCaseCounts query for all data")
		return allCountriesAggregatedData, nil
	}
	log.Printf("GetCountryCaseCounts query from: %s, to: %s, country: %s\n", from, to, country)
	agg, err := aggregateDataBetweenDates(from, to, country)
	return aggregateCountryDataFromStatesAggregate(agg), err
}

// GetWorldCaseCounts : get case counts for the world.
func GetWorldCaseCounts(from string, to string) ([]CaseCount, error) {
	if from == "" && to == "" {
		log.Println("GetWorldCaseCounts query for all data")
		return worldCaseCountsCache, nil
	}
	log.Printf("GetWorldCaseCounts query from: %s, to: %s\n", from, to)
	return getWorldDataBetweenDates(from, to)
}

func setDateBoundariesAndAllAggregatedData(headerRow []string) {
	firstDate, _ = time.Parse(dateformat.CasesDateFormat, headerRow[4])
	lastDate, _ = time.Parse(dateformat.CasesDateFormat, headerRow[len(headerRow)-1])
	worldCaseCountsCache = aggregateWorldData(caseCountsCache)
	allAggregatedData, _ = aggregateDataBetweenDates("", "", "")
	allCountriesAggregatedData = aggregateCountryDataFromStatesAggregate(allAggregatedData)
	countryCaseCountsCache = aggregateCountryDataFromCaseCounts(caseCountsCache)
}
