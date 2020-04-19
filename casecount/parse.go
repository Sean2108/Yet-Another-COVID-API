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

func extractCaseCounts(headerRow []string, confirmedData [][]string, deathsData [][]string) map[string]map[string]CaseCounts {
	caseCountsMap := make(map[string]map[string]CaseCounts)
	numRows := len(confirmedData)
	ch := make(chan extractedInformation, numRows-1)
	wg := sync.WaitGroup{}
	for rowIndex := 1; rowIndex < numRows; rowIndex++ {
		wg.Add(1)
		go getCaseCountsDataForState(headerRow, confirmedData[rowIndex], deathsData[rowIndex], ch, &wg)
	}
	wg.Wait()
	close(ch)
	for item := range ch {
		if _, ok := caseCountsMap[item.country]; !ok {
			caseCountsMap[item.country] = make(map[string]CaseCounts)
		}
		caseCountsMap[item.country][item.state] = item.counts
	}
	return caseCountsMap
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

func getCaseCountsArrayForState(headerRow []string, confirmedRow []string, deathsRow []string) []CaseCount {
	var counts []CaseCount
	for colIndex := 4; colIndex < len(confirmedRow); colIndex++ {
		confirmedCount, confirmedErr := strconv.Atoi(confirmedRow[colIndex])
		if confirmedErr != nil {
			log.Fatal(confirmedErr.Error())
		}
		deathsCount, deathsErr := strconv.Atoi(deathsRow[colIndex])
		if deathsErr != nil {
			log.Fatal(deathsErr.Error())
		}
		caseCountItem := CaseCount{headerRow[colIndex], statistics{confirmedCount, deathsCount}}
		counts = append(counts, caseCountItem)
	}
	return counts
}

func getCaseCountsDataForState(headerRow []string, confirmedRow []string, deathsRow []string, ch chan extractedInformation, wg *sync.WaitGroup) {
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
	ch <- extractedInformation{confirmedRow[0], confirmedRow[1], CaseCounts{Location{float32(lat), float32(long)}, counts}}
	wg.Done()
}
