package casecount

import (
	"log"
	"net/http"
	"sync"
	"time"

	"yet-another-covid-map-api/dateformat"
	"yet-another-covid-map-api/utils"
)

const (
	confirmedURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv"
	deathsURL    = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv"
	recoveredURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_recovered_global.csv"

	usConfirmedURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_US.csv"
	usDeathsURL    = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_US.csv"
)

var (
	// cache the query for getting all data for all states and all countries, because it is the most heavily used
	stateAggregatedMap   map[string]CountryWithStatesAggregated
	countryAggregatedMap map[string]CountryAggregated

	caseCountsMap        map[string]CountryWithStates
	countryCaseCountsMap map[string]Country
	worldCaseCountsCache []CaseCount

	lastDate  time.Time
	firstDate time.Time

	mux sync.Mutex

	client utils.HTTPClient
)

func init() {
	client = &http.Client{}
}

// UpdateCaseCounts : Pull data from the John Hopkins CSV files on GitHub, store the result in a cache and also cache the aggregate data for the entire period
func UpdateCaseCounts() {
	log.Println("Updating case counts")
	confirmedData, deathsData, recoveredData, usConfirmedData, usDeathsData, ok := getData()
	if !ok || len(confirmedData) < 2 || len(confirmedData) != len(deathsData) {
		log.Println("New data is faulty, continuing to use old data.")
		return
	}
	headerRow := confirmedData[0]
	baseCaseCountsMap := extractCaseCounts(headerRow, confirmedData, deathsData, recoveredData)
	usCaseCounts := extractUSCaseCounts(usConfirmedData, usDeathsData)
	mux.Lock()
	defer mux.Unlock()
	caseCountsMap = mergeCaseCountsWithUS(baseCaseCountsMap, usCaseCounts)
	setDateBoundariesAndAllAggregatedData(headerRow)
}

// GetCaseCountsWithDayData : get case counts for states but without aggregating the counts, so a list of days with number of confirmed cases and deaths on each day is returned
func GetCaseCountsWithDayData(from string, to string, country string) (map[string]CountryWithStates, error) {
	if from == "" && to == "" && country == "" {
		log.Println("GetCaseCounts query for all data with per day information")
		return caseCountsMap, nil
	}
	log.Printf("GetCaseCountsWithDayData query from: %s, to: %s, country: %s\n", from, to, country)
	return filterCaseCounts(from, to, country)
}

// GetCountryCaseCountsWithDayData : get case counts for countries but without aggregating the counts, so a list of days with number of confirmed cases and deaths on each day is returned
func GetCountryCaseCountsWithDayData(from string, to string, country string) (map[string]Country, error) {
	if from == "" && to == "" && country == "" {
		log.Println("GetCountryCaseCounts query for all data with per day information")
		return countryCaseCountsMap, nil
	}
	log.Printf("GetCountryCaseCountsWithDayData query from: %s, to: %s, country: %s\n", from, to, country)
	filtered, err := filterCaseCounts(from, to, country)
	return aggregateCountryDataFromCaseCounts(filtered), err
}

// GetCaseCounts : get case counts for all states between from date and to date. Return case counts for entire period if from and to dates are empty strings
func GetCaseCounts(from string, to string, country string) (map[string]CountryWithStatesAggregated, error) {
	if from == "" && to == "" && country == "" {
		log.Println("GetCaseCounts query for all data")
		return stateAggregatedMap, nil
	}
	log.Printf("GetCaseCounts query from: %s, to: %s, country: %s\n", from, to, country)
	return aggregateDataBetweenDates(from, to, country)
}

// GetCountryCaseCounts : get case counts for all countries between from date and to date. Return case counts for entire period if from and to dates are empty strings
func GetCountryCaseCounts(from string, to string, country string) (map[string]CountryAggregated, error) {
	if from == "" && to == "" && country == "" {
		log.Println("GetCountryCaseCounts query for all data")
		return countryAggregatedMap, nil
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
	stateAggregatedMap, _ = aggregateDataBetweenDates("", "", "")
	countryAggregatedMap = aggregateCountryDataFromStatesAggregate(stateAggregatedMap)
	countryCaseCountsMap = aggregateCountryDataFromCaseCounts(caseCountsMap)
	worldCaseCountsCache = aggregateWorldData(countryCaseCountsMap)
}
