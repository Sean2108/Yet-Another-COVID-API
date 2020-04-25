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

func extractCaseCounts(headerRow []string, confirmedData [][]string, deathsData [][]string, recoveredData [][]string) map[string]map[string]CaseCounts {
	caseCountsMap := make(map[string]map[string]CaseCounts)
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
			caseCountsMap[item.country] = make(map[string]CaseCounts)
		}
		caseCountsMap[item.country][item.state] = item.counts
	}
	for item := range recoveredCh {
		if state, ok := caseCountsMap[item.country][item.state]; ok {
			for i, count := range item.counts.Counts {
				state.Counts[i].Recovered = count.Recovered
			}
		} else {
			caseCountsMap[item.country][item.state] = item.counts
		}
	}
	return caseCountsMap
}

func mergeCaseCountsWithUS(caseCountsMap map[string]map[string]CaseCounts, usCaseCounts map[string]CaseCounts) map[string]map[string]CaseCounts {
	usInfo := caseCountsMap["US"]
	for i := range usInfo[""].Counts {
		usInfo[""].Counts[i].Confirmed = 0
		usInfo[""].Counts[i].Deaths = 0
	}
	for key, value := range usCaseCounts {
		usInfo[key] = value
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
				log.Fatal(err.Error())
			}
			long, err := strconv.ParseFloat(confirmedRow[9], 32)
			if err != nil {
				log.Fatal(err.Error())
			}
			if counts, ok := getCaseCountsArray(headerRow, confirmedRow, deathsData[rowIndex], nil, 11, 1); ok {
				usInfo[state] = CaseCounts{Location{float32(lat), float32(long)}, counts}
			}
		}
	}
	return usInfo
}

func getData(url string) [][]string {
	data, err := utils.ReadCSVFromURL(client, url)
	if err != nil {
		log.Fatal(err.Error())
	}
	return data
}

func getColumnValue(row []string, colIndex int) (int, bool) {
	if row != nil {
		count, err := strconv.Atoi(row[colIndex])
		if err != nil {
			log.Fatal(err.Error())
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
	for colIndex := startIndex; colIndex < len(headerRow); colIndex++ {
		confirmedCount, confirmedOk := getColumnValue(confirmedRow, colIndex)
		deathsCount, deathsOk := getColumnValue(deathsRow, colIndex+deathsColOffset)
		recoveredCount, recoveredOk := getColumnValue(recoveredRow, colIndex)
		if !(confirmedOk && deathsOk && recoveredOk) {
			return nil, false
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
		log.Fatal(latError.Error())
	}
	long, longError := strconv.ParseFloat(rowDetails[3], 32)
	if longError != nil {
		log.Fatal(longError.Error())
	}

	ch <- extractedInformation{rowDetails[0], iso, CaseCounts{Location{float32(lat), float32(long)}, counts}}
}
