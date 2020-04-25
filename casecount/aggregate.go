package casecount

import (
	"fmt"
	"strings"
	"sync"
	"yet-another-covid-map-api/utils"
)

type aggregatedCaseCountsMap struct {
	country string
	info    CountryWithStatesAggregated
}

type countryMap struct {
	country string
	info    Country
}

type countryAggMap struct {
	country string
	info    CountryAggregated
}

func copyAndFilterCaseCountsMap(countryInfo CountryWithStates, fromIndex int, toIndex int) CountryWithStates {
	newInfo := CountryWithStates{countryInfo.Name, make(map[string]CaseCounts, len(countryInfo.States))}
	for state, stateInfo := range countryInfo.States {
		newStateInfo := CaseCounts{stateInfo.LocationAndPopulation, stateInfo.Counts[fromIndex : toIndex+1]}
		newInfo.States[state] = newStateInfo
	}
	return newInfo
}

func aggregateCaseCountsMap(countryInfo CountryWithStates, fromIndex int, toIndex int) CountryWithStatesAggregated {
	newInfo := CountryWithStatesAggregated{countryInfo.Name, make(map[string]CaseCountsAggregated, len(countryInfo.States))}
	for state, stateInfo := range countryInfo.States {
		confirmedSum, deathsSum, recoveredSum := getStatisticsSum(stateInfo.Counts, fromIndex, toIndex)
		newStateInfo := CaseCountsAggregated{stateInfo.LocationAndPopulation, statistics{confirmedSum, deathsSum, recoveredSum}}
		newInfo.States[state] = newStateInfo
	}
	return newInfo
}

func filterCaseCounts(from string, to string, country string) (map[string]CountryWithStates, error) {
	fromIndex, toIndex := getFromAndToIndices(from, to)
	filteredCaseCounts := make(map[string]CountryWithStates)
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
		countryName, _ := utils.GetCountryFromAbbreviation(countryKey)
		if country == "" || strings.ToLower(country) == strings.ToLower(countryName) {
			filteredCaseCounts[countryKey] = copyAndFilterCaseCountsMap(countryInfo, fromIndex, toIndex)
		}
	}
	return filteredCaseCounts, nil
}

func syncAggregateCaseCountsMap(countryKey string, countryInfo CountryWithStates, fromIndex int, toIndex int, country string, ch chan aggregatedCaseCountsMap, wg *sync.WaitGroup) {
	countryName, _ := utils.GetCountryFromAbbreviation(countryKey)
	if country == "" || strings.ToLower(countryName) == strings.ToLower(country) {
		info := aggregateCaseCountsMap(countryInfo, fromIndex, toIndex)
		ch <- aggregatedCaseCountsMap{countryKey, info}
	}
	wg.Done()
}

func aggregateDataBetweenDates(from string, to string, country string) (map[string]CountryWithStatesAggregated, error) {
	fromIndex, toIndex := getFromAndToIndices(from, to)
	aggregatedData := make(map[string]CountryWithStatesAggregated)
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
	return aggregatedData, nil
}

func syncSumStates(country string, countryInfo map[string]CaseCounts, ch chan countryMap, wg *sync.WaitGroup) {
	var latSum, longSum float32
	var count, population int
	var counts []CaseCount
	for _, stateInfo := range countryInfo {
		population += stateInfo.Population
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
				counts[index].Recovered += stateInfo.Counts[index].Recovered
			}
		}
	}
	var lat, long float32
	if info, ok := countryInfo[""]; ok {
		lat, long = info.Lat, info.Long
		population = info.Population
	} else {
		countF := float32(count)
		lat, long = latSum/countF, longSum/countF
	}
	countryName, _ := utils.GetCountryFromAbbreviation(country)
	ch <- countryMap{country, Country{countryName, CaseCounts{LocationAndPopulation{lat, long, population}, counts}}}
	wg.Done()
}

func aggregateCountryDataFromCaseCounts(caseCountsMap map[string]CountryWithStates) map[string]Country {
	ch := make(chan countryMap, len(caseCountsMap))
	wg := sync.WaitGroup{}
	aggregatedData := make(map[string]Country)
	for country, countryInfo := range caseCountsMap {
		wg.Add(1)
		go syncSumStates(country, countryInfo.States, ch, &wg)
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
	var count, confirmed, deaths, recovered, population int
	for _, stateInfo := range countryInfo {
		latSum += stateInfo.Lat
		longSum += stateInfo.Long
		count++
		confirmed += stateInfo.Confirmed
		deaths += stateInfo.Deaths
		recovered += stateInfo.Recovered
		population += stateInfo.Population
	}
	var lat, long float32
	if info, ok := countryInfo[""]; ok {
		lat, long = info.Lat, info.Long
		population = info.Population
	} else {
		countF := float32(count)
		lat, long = latSum/countF, longSum/countF
	}
	countryName, _ := utils.GetCountryFromAbbreviation(country)
	ch <- countryAggMap{country, CountryAggregated{countryName, CaseCountsAggregated{LocationAndPopulation{lat, long, population}, statistics{confirmed, deaths, recovered}}}}
	wg.Done()
}

func aggregateCountryDataFromStatesAggregate(caseCountsMap map[string]CountryWithStatesAggregated) map[string]CountryAggregated {
	ch := make(chan countryAggMap, len(caseCountsMap))
	wg := sync.WaitGroup{}
	aggregatedData := make(map[string]CountryAggregated)
	for country, countryInfo := range caseCountsMap {
		wg.Add(1)
		go syncSumStatesAggregated(country, countryInfo.States, ch, &wg)
	}
	wg.Wait()
	close(ch)
	for caseCountsAgg := range ch {
		aggregatedData[caseCountsAgg.country] = caseCountsAgg.info
	}
	return aggregatedData
}

func aggregateWorldData(caseCounts map[string]Country) []CaseCount {
	var counts []CaseCount
	for _, countryInfo := range caseCounts {
		if counts == nil {
			counts = make([]CaseCount, len(countryInfo.Counts))
			copy(counts, countryInfo.Counts)
		} else {
			for index := range counts {
				counts[index].Confirmed += countryInfo.Counts[index].Confirmed
				counts[index].Deaths += countryInfo.Counts[index].Deaths
				counts[index].Recovered += countryInfo.Counts[index].Recovered
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
