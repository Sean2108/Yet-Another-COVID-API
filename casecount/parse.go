package casecount

import (
	"log"
	"strconv"
	"sync"
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

func getData() ([][]string, [][]string, [][]string) {
	confirmedData, confirmedError := readCSVFromURL(confirmedURL)
	deathsData, deathsError := readCSVFromURL(deathsURL)
	recoveredData, recoveredError := readCSVFromURL(recoveredURL)
	if confirmedError != nil {
		log.Fatal(confirmedError.Error())
	}
	if deathsError != nil {
		log.Fatal(deathsError.Error())
	}
	if recoveredError != nil {
		log.Fatal(recoveredError.Error())
	}
	if len(confirmedData) < 2 || len(confirmedData) != len(deathsData) {
		log.Fatal("Invalid CSV files obtained")
	}
	return confirmedData, deathsData, recoveredData
}

func getColumnValue(row []string, colIndex int) int {
	if row != nil {
		count, err := strconv.Atoi(row[colIndex])
		if err != nil {
			log.Fatal(err.Error())
		}
		return count
	}
	return 0
}

func getCaseCountsArray(headerRow []string, confirmedRow []string, deathsRow []string, recoveredRow []string) []CaseCount {
	var counts []CaseCount
	for colIndex := 4; colIndex < len(headerRow); colIndex++ {
		confirmedCount := getColumnValue(confirmedRow, colIndex)
		deathsCount := getColumnValue(deathsRow, colIndex)
		recoveredCount := getColumnValue(recoveredRow, colIndex)
		caseCountItem := CaseCount{headerRow[colIndex], statistics{confirmedCount, deathsCount, recoveredCount}}
		counts = append(counts, caseCountItem)
	}
	return counts
}

func getCaseCountsData(headerRow []string, confirmedRow []string, deathsRow []string, recoveredRow []string, rowDetails []string, ch chan extractedInformation, wg *sync.WaitGroup) {
	// skip faulty entry in data
	if rowDetails[0] == "Diamond Princess" {
		wg.Done()
		return
	}
	counts := getCaseCountsArray(headerRow, confirmedRow, deathsRow, recoveredRow)
	lat, latError := strconv.ParseFloat(rowDetails[2], 32)
	if latError != nil {
		log.Fatal(latError.Error())
	}
	long, longError := strconv.ParseFloat(rowDetails[3], 32)
	if longError != nil {
		log.Fatal(longError.Error())
	}
	ch <- extractedInformation{rowDetails[0], rowDetails[1], CaseCounts{Location{float32(lat), float32(long)}, counts}}
	wg.Done()
}
