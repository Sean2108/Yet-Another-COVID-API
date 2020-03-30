package casecount

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
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

type caseCounts struct {
	stateInformation
	Counts []caseCount
}

// CaseCountsAggregated : contains the information about the state, country and the latitude/longitude as well as the number of confirmed cases/deaths
type CaseCountsAggregated struct {
	stateInformation
	statistics
}

const inputDateFormat = "1/2/06"

var confirmedURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_confirmed_global.csv"
var deathsURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_covid19_deaths_global.csv"

var caseCountsCache []caseCounts

// AllAggregatedData : cache the query for getting all data, because it is the most heavily used
var AllAggregatedData []CaseCountsAggregated

var lastDate time.Time
var firstDate time.Time

// UpdateCaseCounts : Pull data from the John Hopkins CSV files on GitHub, store the result in a cache and also cache the aggregate data for the entire period
func UpdateCaseCounts() {
	log.Println("Updating case counts")
	caseCountsCache = nil
	confirmedData, deathsData := getData()
	headerRow := confirmedData[0]
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
	firstDate, _ = time.Parse(inputDateFormat, headerRow[4])
	lastDate, _ = time.Parse(inputDateFormat, headerRow[len(headerRow)-1])
	AllAggregatedData = aggregateDataBetweenDates("", "")
}

// GetCaseCounts : get case counts for all states between from date and to date. Return case counts for entire period if from and to dates are empty strings
func GetCaseCounts(from string, to string) []CaseCountsAggregated {
	if from == "" && to == "" {
		return AllAggregatedData
	}
	return aggregateDataBetweenDates(from, to)
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
		toIndex = getDaysBetweenDates(toDate, lastDate)
	}
	return fromIndex, toIndex
}

func getStatisticsSum(input []caseCount) (int, int) {
	confirmedSum := 0
	deathsSum := 0
	for _, item := range input {
		confirmedSum += item.Confirmed
		deathsSum += item.Deaths
	}
	return confirmedSum, deathsSum
}

func convertToAggregatedElement(caseCountsItem caseCounts, from int, to int, ch chan CaseCountsAggregated, wg *sync.WaitGroup) {
	confirmedSum, deathsSum := getStatisticsSum(caseCountsItem.Counts[from:to])
	ch <- CaseCountsAggregated{stateInformation{caseCountsItem.State, caseCountsItem.Country, caseCountsItem.Lat, caseCountsItem.Long}, statistics{confirmedSum, deathsSum}}
	wg.Done()
}

func aggregateDataBetweenDates(from string, to string) []CaseCountsAggregated {
	fromIndex, toIndex := getFromAndToIndices(from, to)
	ch := make(chan CaseCountsAggregated, len(caseCountsCache))
	wg := sync.WaitGroup{}
	for _, caseCountsItem := range caseCountsCache {
		wg.Add(1)
		go convertToAggregatedElement(caseCountsItem, fromIndex, toIndex, ch, &wg)
	}
	wg.Wait()
	close(ch)
	var aggregatedData []CaseCountsAggregated
	for caseCountsAgg := range ch {
		aggregatedData = append(aggregatedData, caseCountsAgg)
	}
	return aggregatedData
}
