package casecount

import (
	"log"
	"strconv"
	"sync"
	"yet-another-covid-map-api/utils"
)

type extractedInformation struct {
	state   string
	country string
	counts  CaseCounts
}

func extractCaseCounts(headerRow []string, confirmedData [][]string, deathsData [][]string, recoveredData [][]string) map[string]CountryWithStates {
	caseCountsMap := make(map[string]CountryWithStates)
	numRows := len(confirmedData)
	ch := make(chan extractedInformation, numRows-1)
	recoveredCh := make(chan extractedInformation, len(recoveredData)-1)
	wg := sync.WaitGroup{}
	for rowIndex := 1; rowIndex < numRows; rowIndex++ {
		wg.Add(1)
		go getCaseCountsData(headerRow, confirmedData[rowIndex], deathsData[rowIndex], nil, confirmedData[rowIndex][0:4], ch, &wg)
	}
	for rowIndex := 1; rowIndex < len(recoveredData); rowIndex++ {
		wg.Add(1)
		go getCaseCountsData(headerRow, nil, nil, recoveredData[rowIndex], recoveredData[rowIndex][0:4], recoveredCh, &wg)
	}
	wg.Wait()
	close(ch)
	close(recoveredCh)
	for item := range ch {
		if _, ok := caseCountsMap[item.country]; !ok {
			countryName, _ := utils.GetCountryFromAbbreviation(item.country)
			countryInfo := CountryWithStates{countryName, map[string]CaseCounts{}}
			caseCountsMap[item.country] = countryInfo
		}
		caseCountsMap[item.country].States[item.state] = item.counts
	}
	for item := range recoveredCh {
		if state, ok := caseCountsMap[item.country].States[item.state]; ok {
			for i, count := range item.counts.Counts {
				state.Counts[i].Recovered = count.Recovered
			}
		} else {
			if _, ok := caseCountsMap[item.country]; !ok {
				countryName, _ := utils.GetCountryFromAbbreviation(item.country)
				caseCountsMap[item.country] = CountryWithStates{countryName, map[string]CaseCounts{}}
			}
			caseCountsMap[item.country].States[item.state] = item.counts
		}
	}
	return caseCountsMap
}

func mergeCaseCountsWithUS(caseCountsMap map[string]CountryWithStates, usCaseCounts map[string]CaseCounts) map[string]CountryWithStates {
	usInfo := caseCountsMap["US"]
	for i := range usInfo.States[""].Counts {
		usInfo.States[""].Counts[i].Confirmed = 0
		usInfo.States[""].Counts[i].Deaths = 0
	}
	for key, value := range usCaseCounts {
		usInfo.States[key] = value
	}
	return caseCountsMap
}

func extractUSCaseCounts(confirmedData [][]string, deathsData [][]string) map[string]CaseCounts {
	headerRow := confirmedData[0]
	usInfo := make(map[string]CaseCounts)
	for rowIndex := 1; rowIndex < len(confirmedData); rowIndex++ {
		confirmedRow := confirmedData[rowIndex]
		state := confirmedRow[6]
		if stateInfo, ok := usInfo[state]; ok {
			if counts, ok := getCaseCountsArray(headerRow, confirmedRow, deathsData[rowIndex], nil, 11, 1); ok {
				for i, count := range counts {
					stateInfo.Counts[i].Confirmed += count.Confirmed
					stateInfo.Counts[i].Deaths += count.Deaths
				}
			}
		} else {
			lat, err := strconv.ParseFloat(confirmedRow[8], 32)
			if err != nil {
				log.Println(err.Error())
				return nil
			}
			long, err := strconv.ParseFloat(confirmedRow[9], 32)
			if err != nil {
				log.Println(err.Error())
				return nil
			}
			if counts, ok := getCaseCountsArray(headerRow, confirmedRow, deathsData[rowIndex], nil, 11, 1); ok {
				usInfo[state] = CaseCounts{LocationAndPopulation{float32(lat), float32(long), utils.StatePopulationLookup["US"][state]}, counts}
			}
		}
	}
	return usInfo
}

func getData() ([][]string, [][]string, [][]string, [][]string, [][]string, bool) {
	confirmedData, confirmedOk := utils.ReadCSVFromURL(client, confirmedURL)
	deathsData, deathsOk := utils.ReadCSVFromURL(client, deathsURL)
	recoveredData, recoveredOk := utils.ReadCSVFromURL(client, recoveredURL)
	usConfirmedData, usConfirmedOk := utils.ReadCSVFromURL(client, usConfirmedURL)
	usDeathsData, usDeathsOk := utils.ReadCSVFromURL(client, usDeathsURL)
	return confirmedData, deathsData, recoveredData, usConfirmedData, usDeathsData,
		confirmedOk && deathsOk && recoveredOk && usConfirmedOk && usDeathsOk
}

func getColumnValue(row []string, colIndex int) (int, bool) {
	if row != nil {
		count, err := strconv.Atoi(row[colIndex])
		if err != nil {
			log.Println(err.Error())
			return 0, false
		}
		if count < 0 {
			return 0, false
		}
		return count, true
	}
	return 0, true
}

func getCaseCountsArray(headerRow []string, confirmedRow []string, deathsRow []string, recoveredRow []string, startIndex int, deathsColOffset int) ([]CaseCount, bool) {
	var counts []CaseCount
	previousRecovered := 0
	for colIndex := startIndex; colIndex < len(headerRow); colIndex++ {
		confirmedCount, confirmedOk := getColumnValue(confirmedRow, colIndex)
		deathsCount, deathsOk := getColumnValue(deathsRow, colIndex+deathsColOffset)
		recoveredCount, recoveredOk := getColumnValue(recoveredRow, colIndex)
		if !(confirmedOk && deathsOk && recoveredOk) {
			return nil, false
		}
		if recoveredCount == 0 {
			// workaround for https://github.com/CSSEGISandData/COVID-19/issues/4465,
			// recovery data is discontinued
			recoveredCount = previousRecovered
		} else {
			previousRecovered = recoveredCount
		}
		caseCountItem := CaseCount{headerRow[colIndex], statistics{confirmedCount, deathsCount, recoveredCount}}
		counts = append(counts, caseCountItem)
	}
	return counts, true
}

func getCaseCountsData(headerRow []string, confirmedRow []string, deathsRow []string, recoveredRow []string, rowDetails []string, ch chan extractedInformation, wg *sync.WaitGroup) {
	defer wg.Done()
	counts, ok := getCaseCountsArray(headerRow, confirmedRow, deathsRow, recoveredRow, 4, 0)
	iso, lookupOk := utils.GetAbbreviationFromCountry(rowDetails[1])
	if !ok || !lookupOk {
		return
	}
	lat, latError := strconv.ParseFloat(rowDetails[2], 32)
	if latError != nil {
		log.Println(latError.Error())
		return
	}
	long, longError := strconv.ParseFloat(rowDetails[3], 32)
	if longError != nil {
		log.Println(longError.Error())
		return
	}
	ch <- extractedInformation{rowDetails[0], iso, CaseCounts{LocationAndPopulation{float32(lat), float32(long), utils.StatePopulationLookup[iso][rowDetails[0]]}, counts}}
}
