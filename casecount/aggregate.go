package casecount

import (
	"fmt"
	"strings"
	"sync"
)

type aggregatedCaseCountsMap struct {
	country string
	info    map[string]CaseCountsAggregated
}

type countryMap struct {
	country string
	info    CaseCounts
}

type countryAggMap struct {
	country string
	info    CaseCountsAggregated
}

func copyAndFilterCaseCountsMap(countryInfo map[string]CaseCounts, fromIndex int, toIndex int) map[string]CaseCounts {
	newInfo := make(map[string]CaseCounts, len(countryInfo))
	for state, stateInfo := range countryInfo {
		newStateInfo := CaseCounts{stateInfo.Location, stateInfo.Counts[fromIndex : toIndex+1]}
		newInfo[state] = newStateInfo
	}
	return newInfo
}

func aggregateCaseCountsMap(countryInfo map[string]CaseCounts, fromIndex int, toIndex int) map[string]CaseCountsAggregated {
	newInfo := make(map[string]CaseCountsAggregated, len(countryInfo))
	for state, stateInfo := range countryInfo {
		confirmedSum, deathsSum := getStatisticsSum(stateInfo.Counts, fromIndex, toIndex)
		newStateInfo := CaseCountsAggregated{stateInfo.Location, statistics{confirmedSum, deathsSum}}
		newInfo[state] = newStateInfo
	}
	return newInfo
}

func filterCaseCounts(from string, to string, country string) (map[string]map[string]CaseCounts, error) {
	fromIndex, toIndex := getFromAndToIndices(from, to)
	filteredCaseCounts := make(map[string]map[string]CaseCounts)
	if fromIndex > toIndex {
		return filteredCaseCounts, fmt.Errorf("From date %s cannot be after to date %s", from, to)
	}
	if country != "" {
		if countryInfo, ok := caseCountsMap[country]; ok {
			filteredCaseCounts[country] = copyAndFilterCaseCountsMap(countryInfo, fromIndex, toIndex)
			return filteredCaseCounts, nil
		}
	}
	for countryKey, countryInfo := range caseCountsMap {
		if country == "" || strings.ToLower(country) == strings.ToLower(countryKey) {
			filteredCaseCounts[countryKey] = copyAndFilterCaseCountsMap(countryInfo, fromIndex, toIndex)
		}
	}
	var err error
	if country != "" && len(filteredCaseCounts) == 0 {
		err = fmt.Errorf("Country %s not found, did you mean: %s?", country, findClosestMatchToCountryName(country))
	}
	return filteredCaseCounts, err
}

func syncAggregateCaseCountsMap(countryKey string, countryInfo map[string]CaseCounts, fromIndex int, toIndex int, country string, ch chan aggregatedCaseCountsMap, wg *sync.WaitGroup) {
	if country == "" || strings.ToLower(countryKey) == strings.ToLower(country) {
		info := aggregateCaseCountsMap(countryInfo, fromIndex, toIndex)
		ch <- aggregatedCaseCountsMap{countryKey, info}
	}
	wg.Done()
}

func aggregateDataBetweenDates(from string, to string, country string) (map[string]map[string]CaseCountsAggregated, error) {
	fromIndex, toIndex := getFromAndToIndices(from, to)
	aggregatedData := make(map[string]map[string]CaseCountsAggregated)
	if fromIndex > toIndex {
		return aggregatedData, fmt.Errorf("From date %s cannot be after to date %s", from, to)
	}
	if country != "" {
		if countryInfo, ok := caseCountsMap[country]; ok {
			aggregatedData[country] = aggregateCaseCountsMap(countryInfo, fromIndex, toIndex)
			return aggregatedData, nil
		}
	}
	ch := make(chan aggregatedCaseCountsMap, len(caseCountsMap))
	wg := sync.WaitGroup{}
	for countryKey, countryInfo := range caseCountsMap {
		wg.Add(1)
		go syncAggregateCaseCountsMap(countryKey, countryInfo, fromIndex, toIndex, country, ch, &wg)
	}
	wg.Wait()
	close(ch)
	for caseCountsAgg := range ch {
		aggregatedData[caseCountsAgg.country] = caseCountsAgg.info
	}
	var err error
	if country != "" && len(aggregatedData) == 0 {
		err = fmt.Errorf("Country %s not found, did you mean: %s?", country, findClosestMatchToCountryName(country))
	}
	return aggregatedData, err
}

func syncSumStates(country string, countryInfo map[string]CaseCounts, ch chan countryMap, wg *sync.WaitGroup) {
	var latSum, longSum float32
	count := 0
	var counts []CaseCount
	for _, stateInfo := range countryInfo {
		latSum += stateInfo.Lat
		longSum += stateInfo.Long
		count++
		if counts == nil {
			counts = make([]CaseCount, len(stateInfo.Counts))
			copy(counts, stateInfo.Counts)
		} else {
			for index := range counts {
				counts[index].Confirmed += stateInfo.Counts[index].Confirmed
				counts[index].Deaths += stateInfo.Counts[index].Deaths
			}
		}
	}
	countF := float32(count)
	ch <- countryMap{country, CaseCounts{Location{latSum / countF, longSum / countF}, counts}}
	wg.Done()
}

func aggregateCountryDataFromCaseCounts(caseCountsMap map[string]map[string]CaseCounts) map[string]CaseCounts {
	ch := make(chan countryMap, len(caseCountsMap))
	wg := sync.WaitGroup{}
	aggregatedData := make(map[string]CaseCounts)
	for country, countryInfo := range caseCountsMap {
		wg.Add(1)
		go syncSumStates(country, countryInfo, ch, &wg)
	}
	wg.Wait()
	close(ch)
	for caseCountsAgg := range ch {
		aggregatedData[caseCountsAgg.country] = caseCountsAgg.info
	}
	return aggregatedData
}

func syncSumStatesAggregated(country string, countryInfo map[string]CaseCountsAggregated, ch chan countryAggMap, wg *sync.WaitGroup) {
	var latSum, longSum float32
	count := 0
	confirmed := 0
	deaths := 0
	for _, stateInfo := range countryInfo {
		latSum += stateInfo.Lat
		longSum += stateInfo.Long
		count++
		confirmed += stateInfo.Confirmed
		deaths += stateInfo.Deaths
	}
	countF := float32(count)
	ch <- countryAggMap{country, CaseCountsAggregated{Location{latSum / countF, longSum / countF}, statistics{confirmed, deaths}}}
	wg.Done()
}

func aggregateCountryDataFromStatesAggregate(caseCountsMap map[string]map[string]CaseCountsAggregated) map[string]CaseCountsAggregated {
	ch := make(chan countryAggMap, len(caseCountsMap))
	wg := sync.WaitGroup{}
	aggregatedData := make(map[string]CaseCountsAggregated)
	for country, countryInfo := range caseCountsMap {
		wg.Add(1)
		go syncSumStatesAggregated(country, countryInfo, ch, &wg)
	}
	wg.Wait()
	close(ch)
	for caseCountsAgg := range ch {
		aggregatedData[caseCountsAgg.country] = caseCountsAgg.info
	}
	return aggregatedData
}

func aggregateWorldData(caseCounts map[string]map[string]CaseCounts) []CaseCount {
	var counts []CaseCount
	for _, countryInfo := range caseCounts {
		for _, stateInfo := range countryInfo {
			if counts == nil {
				counts = make([]CaseCount, len(stateInfo.Counts))
				copy(counts, stateInfo.Counts)
			} else {
				for index := range counts {
					counts[index].Confirmed += stateInfo.Counts[index].Confirmed
					counts[index].Deaths += stateInfo.Counts[index].Deaths
				}
			}
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
