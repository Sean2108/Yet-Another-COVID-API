package casecount

import (
	"encoding/csv"
	"time"
)

func readCSVFromURL(url string) ([][]string, error) {
	resp, err := client.Get(url)
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

func getDaysBetweenDates(startDate time.Time, endDate time.Time) int {
	return int(endDate.Sub(startDate).Hours() / 24)
}

func getStatisticsSum(input []caseCount, fromIndex int, toIndex int) (int, int) {
	confirmedAtStartDate := 0
	deathsAtStartDate := 0

	if fromIndex >= len(input) || toIndex < 0 {
		return 0, 0
	}

	if fromIndex > 0 {
		confirmedAtStartDate = input[fromIndex-1].Confirmed
		deathsAtStartDate = input[fromIndex-1].Deaths
	}
	if toIndex >= len(input) {
		toIndex = len(input) - 1
	}

	return input[toIndex].Confirmed - confirmedAtStartDate, input[toIndex].Deaths - deathsAtStartDate
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
