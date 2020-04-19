package casecount

import (
	"fmt"
	"strings"
	"sync"
)

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
	type countryAggregationInformation struct {
		LatSum, LongSum                float32
		ConfirmedSum, DeathsSum, Count int
	}
	var countryInformationMap map[string]countryAggregationInformation
	countryInformationMap = make(map[string]countryAggregationInformation)
	for _, caseCountsAgg := range aggregateDataWithStates {
		if val, ok := countryInformationMap[caseCountsAgg.Country]; ok {
			countryInformationMap[caseCountsAgg.Country] = countryAggregationInformation{val.LatSum + caseCountsAgg.Lat, val.LongSum + caseCountsAgg.Long, val.ConfirmedSum + caseCountsAgg.Confirmed, val.DeathsSum + caseCountsAgg.Deaths, val.Count + 1}
		} else {
			countryInformationMap[caseCountsAgg.Country] = countryAggregationInformation{caseCountsAgg.Lat, caseCountsAgg.Long, caseCountsAgg.Confirmed, caseCountsAgg.Deaths, 1}
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
	type countryAggregationInformation struct {
		LatSum, LongSum float32
		Counts          []CaseCount
		Count           int
	}
	var countryInformationMap map[string]countryAggregationInformation
	countryInformationMap = make(map[string]countryAggregationInformation)
	for _, caseCountsAgg := range caseCounts {
		if val, ok := countryInformationMap[caseCountsAgg.Country]; ok {
			var counts = make([]CaseCount, len(val.Counts))
			copy(counts, val.Counts)
			for index := range counts {
				counts[index].Confirmed += caseCountsAgg.Counts[index].Confirmed
				counts[index].Deaths += caseCountsAgg.Counts[index].Deaths
			}
			countryInformationMap[caseCountsAgg.Country] = countryAggregationInformation{val.LatSum + caseCountsAgg.Lat, val.LongSum + caseCountsAgg.Long, counts, val.Count + 1}
		} else {
			countryInformationMap[caseCountsAgg.Country] = countryAggregationInformation{caseCountsAgg.Lat, caseCountsAgg.Long, caseCountsAgg.Counts, 1}
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

func aggregateWorldData(caseCounts []CaseCounts) []CaseCount {
	var counts []CaseCount
	counts = make([]CaseCount, len(caseCounts[0].Counts))
	copy(counts, caseCounts[0].Counts)
	for i := 1; i < len(caseCounts); i++ {
		for j, count := range caseCounts[i].Counts {
			counts[j].Confirmed += count.Confirmed
			counts[j].Deaths += count.Deaths
		}
	}
	return counts
}

func getWorldDataBetweenDates(from string, to string) ([]CaseCount, error) {
	fromIndex, toIndex := getFromAndToIndices(from, to)
	if fromIndex > toIndex {
		return nil, fmt.Errorf("From date %s cannot be after to date %s", from, to)
	}
	return worldCaseCountsCache[fromIndex : toIndex+1], nil
}
