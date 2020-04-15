package casecount

import (
	"log"
	"strconv"
	"sync"
)

func extractCaseCounts(headerRow []string, confirmedData [][]string, deathsData [][]string) {
	numRows := len(confirmedData)
	ch := make(chan CaseCounts, numRows-1)
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

func getCaseCountsDataForState(headerRow []string, confirmedRow []string, deathsRow []string, ch chan CaseCounts, wg *sync.WaitGroup) {
	// skip faulty entry in data
	if confirmedRow[0] == "Diamond Princess" {
		wg.Done()
		return
	}
	counts := getCaseCountsArrayForState(headerRow, confirmedRow, deathsRow)
	lat, latError := strconv.ParseFloat(confirmedRow[2], 32)
	if latError != nil {
		log.Fatal(latError.Error())
	}
	long, longError := strconv.ParseFloat(confirmedRow[3], 32)
	if longError != nil {
		log.Fatal(longError.Error())
	}
	ch <- CaseCounts{stateInformation{confirmedRow[0], confirmedRow[1], float32(lat), float32(long)}, counts}
	wg.Done()
}
