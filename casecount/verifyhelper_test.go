package casecount

import (
	"sort"
	"testing"
)

func verifyResultsCaseCountsArr(result []CaseCounts, expectedData []CaseCounts, t *testing.T) {
	sort.Sort(ByCountryAndStateForCaseCounts(result))
	for i, item := range result {
		if !item.equals(expectedData[i]) {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", item, expectedData[i])
		}
	}
}

func verifyResultsCountryCaseCountsArr(result []CountryCaseCounts, expectedData []CountryCaseCounts, t *testing.T) {
	sort.Sort(ByCountryForCaseCounts(result))
	for i, item := range result {
		if !item.equals(expectedData[i]) {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", item, expectedData[i])
		}
	}
}

func verifyResultsCaseCountsAgg(result []CaseCountsAggregated, expectedData []CaseCountsAggregated, t *testing.T) {
	sort.Sort(ByCountryAndStateAgg(result))
	for i, item := range result {
		if item != expectedData[i] {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", item, expectedData[i])
		}
	}
}

func verifyResultsCountryCaseCountsAgg(result []CountryCaseCountsAggregated, expectedData []CountryCaseCountsAggregated, t *testing.T) {
	sort.Sort(ByCountryAgg(result))
	for i, item := range result {
		if item.Country != expectedData[i].Country || int(item.Lat) != int(expectedData[i].Lat) || int(item.Long) != int(expectedData[i].Long) ||
			item.Confirmed != expectedData[i].Confirmed || item.Deaths != expectedData[i].Deaths {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", item, expectedData[i])
		}
	}
}
