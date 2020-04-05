package casecount

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"yet-another-covid-map-api/dateformat"
	"yet-another-covid-map-api/utils"
)

type statistics struct {
	Confirmed int
	Deaths    int
}

type caseCount struct {
	Date string
	statistics
}

type stateInformation struct {
	State   string
	Country string
	Lat     float32
	Long    float32
}

type countryInformation struct {
	Country string
	Lat     float32
	Long    float32
}

// CaseCounts : contains information about the state,country and latitude longitude as well as the per day cumulative number of confirmed cases/deaths
type CaseCounts struct {
	stateInformation
	Counts []caseCount
}

// CountryCaseCounts : contains information about the state,country and latitude longitude as well as the per day cumulative number of confirmed cases/deaths
type CountryCaseCounts struct {
	countryInformation
	Counts []caseCount
}

// CaseCountsAggregated : contains the information about the state, country and the latitude/longitude as well as the number of confirmed cases/deaths
type CaseCountsAggregated struct {
	stateInformation
	statistics
}

// CountryCaseCountsAggregated : contains the information about the country and the latitude/longitude as well as the number of confirmed cases/deaths
type CountryCaseCountsAggregated struct {
	countryInformation
	statistics
}

const (
	confirmedURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv"
	deathsURL    = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv"
)

var (
	// cache the query for getting all data for all states and all countries, because it is the most heavily used
	allAggregatedData          []CaseCountsAggregated
	allCountriesAggregatedData []CountryCaseCountsAggregated

	caseCountsCache        []CaseCounts
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

func setDateBoundariesAndAllAggregatedData(headerRow []string) {
	firstDate, _ = time.Parse(dateformat.CasesDateFormat, headerRow[4])
	lastDate, _ = time.Parse(dateformat.CasesDateFormat, headerRow[len(headerRow)-1])
	allAggregatedData, _ = aggregateDataBetweenDates("", "", "")
	allCountriesAggregatedData = aggregateCountryDataFromStatesAggregate(allAggregatedData)
	countryCaseCountsCache = aggregateCountryDataFromCaseCounts(caseCountsCache)
}

func extractCaseCounts(headerRow []string, confirmedData [][]string, deathsData [][]string) {
	numRows := len(confirmedData)
	ch := make(chan CaseCounts, numRows-1)
	wg := sync.WaitGroup{}
	for rowIndex := 1; rowIndex < numRows; rowIndex++ {
		wg.Add(1)
		go getCaseCountsDataForState(headerRow, confirmedData[rowIndex], deathsData[rowIndex], ch, &wg)
	}
	wg.Wait()
	close(ch)
	for caseCountsItem := range ch {
		caseCountsCache = append(caseCountsCache, caseCountsItem)
	}
}

func getData() ([][]string, [][]string) {
	confirmedData, confirmedError := readCSVFromURL(confirmedURL)
	deathsData, deathsError := readCSVFromURL(deathsURL)
	if confirmedError != nil {
		log.Fatal(confirmedError.Error())
	}
	if deathsError != nil {
		log.Fatal(deathsError.Error())
	}
	if len(confirmedData) < 2 || len(confirmedData) != len(deathsData) {
		log.Fatal("Invalid CSV files obtained")
	}
	return confirmedData, deathsData
}

func getCaseCountsArrayForState(headerRow []string, confirmedRow []string, deathsRow []string) []caseCount {
	var counts []caseCount
	for colIndex := 4; colIndex < len(confirmedRow); colIndex++ {
		confirmedCount, confirmedErr := strconv.Atoi(confirmedRow[colIndex])
		if confirmedErr != nil {
			log.Fatal(confirmedErr.Error())
		}
		deathsCount, deathsErr := strconv.Atoi(deathsRow[colIndex])
		if deathsErr != nil {
			log.Fatal(deathsErr.Error())
		}
		caseCountItem := caseCount{headerRow[colIndex], statistics{confirmedCount, deathsCount}}
		counts = append(counts, caseCountItem)
	}
	return counts
}

func getCaseCountsDataForState(headerRow []string, confirmedRow []string, deathsRow []string, ch chan CaseCounts, wg *sync.WaitGroup) {
	counts := getCaseCountsArrayForState(headerRow, confirmedRow, deathsRow)
	lat, latError := strconv.ParseFloat(confirmedRow[2], 32)
	if latError != nil {
		log.Fatal(latError.Error())
	}
	long, longError := strconv.ParseFloat(confirmedRow[3], 32)
	if longError != nil {
		log.Fatal(longError.Error())
	}
	ch <- CaseCounts{stateInformation{confirmedRow[0], confirmedRow[1], float32(lat), float32(long)}, counts}
	wg.Done()
}

func filterCaseCounts(from string, to string, country string) ([]CaseCounts, error) {
	fromIndex, toIndex := getFromAndToIndices(from, to)
	var filteredCaseCounts []CaseCounts
	if fromIndex > toIndex {
		return filteredCaseCounts, fmt.Errorf("From date %s cannot be after to date %s", from, to)
	}
	for _, caseCountsItem := range caseCountsCache {
		if country == "" || strings.ToLower(caseCountsItem.Country) == strings.ToLower(country) {
			newCaseCountsItem := CaseCounts{caseCountsItem.stateInformation, caseCountsItem.Counts[fromIndex : toIndex+1]}
			filteredCaseCounts = append(filteredCaseCounts, newCaseCountsItem)
		}
	}
	var err error
	if country != "" && len(filteredCaseCounts) == 0 {
		err = fmt.Errorf("Country %s not found, did you mean: %s?", country, findClosestMatchToCountryName(country))
	}
	return filteredCaseCounts, err
}

func convertToAggregatedElement(caseCountsItem CaseCounts, from int, to int, country string, ch chan CaseCountsAggregated, wg *sync.WaitGroup) {
	if country == "" || strings.ToLower(caseCountsItem.Country) == strings.ToLower(country) {
		confirmedSum, deathsSum := getStatisticsSum(caseCountsItem.Counts, from, to)
		ch <- CaseCountsAggregated{stateInformation{caseCountsItem.State, caseCountsItem.Country, caseCountsItem.Lat, caseCountsItem.Long}, statistics{confirmedSum, deathsSum}}
	}
	wg.Done()
}

func aggregateDataBetweenDates(from string, to string, country string) ([]CaseCountsAggregated, error) {
	fromIndex, toIndex := getFromAndToIndices(from, to)
	var aggregatedData []CaseCountsAggregated
	if fromIndex > toIndex {
		return aggregatedData, fmt.Errorf("From date %s cannot be after to date %s", from, to)
	}
	ch := make(chan CaseCountsAggregated, len(caseCountsCache))
	wg := sync.WaitGroup{}
	for _, caseCountsItem := range caseCountsCache {
		wg.Add(1)
		go convertToAggregatedElement(caseCountsItem, fromIndex, toIndex, country, ch, &wg)
	}
	wg.Wait()
	close(ch)
	for caseCountsAgg := range ch {
		aggregatedData = append(aggregatedData, caseCountsAgg)
	}
	var err error
	if country != "" && len(aggregatedData) == 0 {
		err = fmt.Errorf("Country %s not found, did you mean: %s?", country, findClosestMatchToCountryName(country))
	}
	return aggregatedData, err
}

func aggregateCountryDataFromStatesAggregate(aggregateDataWithStates []CaseCountsAggregated) []CountryCaseCountsAggregated {
	type CountryAggregationInformation struct {
		LatSum, LongSum                float32
		ConfirmedSum, DeathsSum, Count int
	}
	var countryInformationMap map[string]CountryAggregationInformation
	countryInformationMap = make(map[string]CountryAggregationInformation)
	for _, caseCountsAgg := range aggregateDataWithStates {
		if val, ok := countryInformationMap[caseCountsAgg.Country]; ok {
			countryInformationMap[caseCountsAgg.Country] = CountryAggregationInformation{val.LatSum + caseCountsAgg.Lat, val.LongSum + caseCountsAgg.Long, val.ConfirmedSum + caseCountsAgg.Confirmed, val.DeathsSum + caseCountsAgg.Deaths, val.Count + 1}
		} else {
			countryInformationMap[caseCountsAgg.Country] = CountryAggregationInformation{caseCountsAgg.Lat, caseCountsAgg.Long, caseCountsAgg.Confirmed, caseCountsAgg.Deaths, 1}
		}
	}
	var aggregatedData []CountryCaseCountsAggregated
	for country, information := range countryInformationMap {
		countF := float32(information.Count)
		countryCaseCountAgg := CountryCaseCountsAggregated{countryInformation{country, information.LatSum / countF, information.LongSum / countF}, statistics{information.ConfirmedSum, information.DeathsSum}}
		aggregatedData = append(aggregatedData, countryCaseCountAgg)
	}
	return aggregatedData
}

func aggregateCountryDataFromCaseCounts(caseCounts []CaseCounts) []CountryCaseCounts {
	type CountryAggregationInformation struct {
		LatSum, LongSum float32
		Counts          []caseCount
		Count           int
	}
	var countryInformationMap map[string]CountryAggregationInformation
	countryInformationMap = make(map[string]CountryAggregationInformation)
	for _, caseCountsAgg := range caseCounts {
		if val, ok := countryInformationMap[caseCountsAgg.Country]; ok {
			var counts = make([]caseCount, len(val.Counts))
			copy(counts, val.Counts)
			for index := range counts {
				counts[index].Confirmed += caseCountsAgg.Counts[index].Confirmed
				counts[index].Deaths += caseCountsAgg.Counts[index].Deaths
			}
			countryInformationMap[caseCountsAgg.Country] = CountryAggregationInformation{val.LatSum + caseCountsAgg.Lat, val.LongSum + caseCountsAgg.Long, counts, val.Count + 1}
		} else {
			countryInformationMap[caseCountsAgg.Country] = CountryAggregationInformation{caseCountsAgg.Lat, caseCountsAgg.Long, caseCountsAgg.Counts, 1}
		}
	}
	var aggregatedData []CountryCaseCounts
	for country, information := range countryInformationMap {
		countF := float32(information.Count)
		countryCaseCountAgg := CountryCaseCounts{countryInformation{country, information.LatSum / countF, information.LongSum / countF}, information.Counts}
		aggregatedData = append(aggregatedData, countryCaseCountAgg)
	}
	return aggregatedData
}
