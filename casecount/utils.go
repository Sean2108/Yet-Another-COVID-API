package casecount

import (
	"time"

	"yet-another-covid-map-api/dateformat"
)

func getDaysBetweenDates(startDate time.Time, endDate time.Time) int {
	return int(endDate.Sub(startDate).Hours() / 24)
}

func getStatisticsSum(input []CaseCount, fromIndex int, toIndex int) (int, int, int) {
	confirmedAtStartDate := 0
	deathsAtStartDate := 0
	recoveredAtStartDate := 0

	if fromIndex >= len(input) || toIndex < 0 {
		return 0, 0, 0
	}

	if fromIndex > 0 {
		confirmedAtStartDate = input[fromIndex-1].Confirmed
		deathsAtStartDate = input[fromIndex-1].Deaths
		recoveredAtStartDate = input[fromIndex-1].Recovered
	}
	if toIndex >= len(input) {
		toIndex = len(input) - 1
	}

	return input[toIndex].Confirmed - confirmedAtStartDate, input[toIndex].Deaths - deathsAtStartDate, input[toIndex].Recovered - recoveredAtStartDate
}

func getFromAndToIndices(from string, to string) (int, int) {
	fromIndex := 0
	toIndex := getDaysBetweenDates(firstDate, lastDate)
	if from == "" && to == "" {
		return fromIndex, toIndex
	}
	fromDate, fromError := time.Parse(dateformat.CasesDateFormat, from)
	toDate, toError := time.Parse(dateformat.CasesDateFormat, to)
	if fromError == nil && fromDate.After(firstDate) {
		fromIndex = getDaysBetweenDates(firstDate, fromDate)
	}
	if toError == nil && toDate.Before(lastDate) {
		toIndex = getDaysBetweenDates(firstDate, toDate)
	}
	return fromIndex, toIndex
}
