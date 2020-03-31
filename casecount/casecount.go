package casecount

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
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

type caseCounts struct {
	stateInformation
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

const inputDateFormat = "1/2/06"

var confirmedURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv"
var deathsURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv"

var caseCountsCache []caseCounts

// cache the query for getting all data for all states and all countries, because it is the most heavily used
var allAggregatedData []CaseCountsAggregated
var allCountriesAggregatedData []CountryCaseCountsAggregated

var lastDate time.Time
var firstDate time.Time

// UpdateCaseCounts : Pull data from the John Hopkins CSV files on GitHub, store the result in a cache and also cache the aggregate data for the entire period
func UpdateCaseCounts() {
	log.Println("Updating case counts")
	caseCountsCache = nil
	confirmedData, deathsData := getData()
	headerRow := confirmedData[0]
	extractCaseCounts(headerRow, confirmedData, deathsData)
	setDateBoundariesAndAllAggregatedData(headerRow)
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
	firstDate, _ = time.Parse(inputDateFormat, headerRow[4])
	lastDate, _ = time.Parse(inputDateFormat, headerRow[len(headerRow)-1])
	allAggregatedData, _ = aggregateDataBetweenDates("", "", "")
	allCountriesAggregatedData = aggregateCountryDataFromStatesAggregate(allAggregatedData)
}

func extractCaseCounts(headerRow []string, confirmedData [][]string, deathsData [][]string) {
	numRows := len(confirmedData)
	ch := make(chan caseCounts, numRows-1)
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

func readCSVFromURL(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
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

func getCaseCountsDataForState(headerRow []string, confirmedRow []string, deathsRow []string, ch chan caseCounts, wg *sync.WaitGroup) {
	counts := getCaseCountsArrayForState(headerRow, confirmedRow, deathsRow)
	lat, latError := strconv.ParseFloat(confirmedRow[2], 32)
	if latError != nil {
		log.Fatal(latError.Error())
	}
	long, longError := strconv.ParseFloat(confirmedRow[3], 32)
	if longError != nil {
		log.Fatal(longError.Error())
	}
	ch <- caseCounts{stateInformation{confirmedRow[0], confirmedRow[1], float32(lat), float32(long)}, counts}
	wg.Done()
}

func getDaysBetweenDates(startDate time.Time, endDate time.Time) int {
	return int(endDate.Sub(startDate).Hours() / 24)
}

func getFromAndToIndices(from string, to string) (int, int) {
	fromIndex := 0
	toIndex := getDaysBetweenDates(firstDate, lastDate)
	if from == "" && to == "" {
		return fromIndex, toIndex
	}
	fromDate, fromError := time.Parse(inputDateFormat, from)
	toDate, toError := time.Parse(inputDateFormat, to)
	if fromError == nil && fromDate.After(firstDate) {
		fromIndex = getDaysBetweenDates(firstDate, fromDate)
	}
	if toError == nil && toDate.Before(lastDate) {
		toIndex = getDaysBetweenDates(firstDate, toDate)
	}
	return fromIndex, toIndex
}

func getStatisticsSum(input []caseCount, fromIndex int, toIndex int) (int, int) {
	confirmedAtStartDate := 0
	deathsAtStartDate := 0
	if fromIndex > 0 {
		confirmedAtStartDate = input[fromIndex-1].Confirmed
		deathsAtStartDate = input[fromIndex-1].Deaths
	}
	return input[toIndex].Confirmed - confirmedAtStartDate, input[toIndex].Deaths - deathsAtStartDate
}

func convertToAggregatedElement(caseCountsItem caseCounts, from int, to int, country string, ch chan CaseCountsAggregated, wg *sync.WaitGroup) {
	if country == "" || strings.ToLower(caseCountsItem.Country) == strings.ToLower(country) {
		confirmedSum, deathsSum := getStatisticsSum(caseCountsItem.Counts, from, to)
		ch <- CaseCountsAggregated{stateInformation{caseCountsItem.State, caseCountsItem.Country, caseCountsItem.Lat, caseCountsItem.Long}, statistics{confirmedSum, deathsSum}}
	}
	wg.Done()
}

func findClosestMatchToCountryName(country string) string {
	minEditDistance := editDistance([]rune(country), []rune(allCountriesAggregatedData[0].Country))
	closestMatch := allCountriesAggregatedData[0].Country
	for _, countryAgg := range allCountriesAggregatedData {
		if editDistance := editDistance([]rune(country), []rune(countryAgg.Country)); editDistance < minEditDistance {
			minEditDistance = editDistance
			closestMatch = countryAgg.Country
		}
	}
	return closestMatch
}

func aggregateDataBetweenDates(from string, to string, country string) ([]CaseCountsAggregated, error) {
	fromIndex, toIndex := getFromAndToIndices(from, to)
	ch := make(chan CaseCountsAggregated, len(caseCountsCache))
	wg := sync.WaitGroup{}
	for _, caseCountsItem := range caseCountsCache {
		wg.Add(1)
		go convertToAggregatedElement(caseCountsItem, fromIndex, toIndex, country, ch, &wg)
	}
	wg.Wait()
	close(ch)
	var aggregatedData []CaseCountsAggregated
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
