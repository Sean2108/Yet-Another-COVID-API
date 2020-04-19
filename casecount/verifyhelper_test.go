package casecount

import (
	"sort"
	"testing"
)

func (a *CaseCounts) equals(b CaseCounts) bool {
	if a.State != b.State || a.Country != b.Country || a.Lat != b.Lat || a.Long != b.Long {
		return false
	}
	if len(a.Counts) != len(b.Counts) {
		return false
	}
	for i, item := range a.Counts {
		if item != b.Counts[i] {
			return false
		}
	}
	return true
}

func (a *CountryCaseCounts) equals(b CountryCaseCounts) bool {
	if a.Country != b.Country || int(a.Lat) != int(b.Lat) || int(a.Long) != int(b.Long) {
		return false
	}
	if len(a.Counts) != len(b.Counts) {
		return false
	}
	for i, item := range a.Counts {
		if item != b.Counts[i] {
			return false
		}
	}
	return true
}

func verifyResultsCaseCountArr(result []CaseCount, expectedData []CaseCount, t *testing.T) {
	for i, item := range result {
		if item != expectedData[i] {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", item, expectedData[i])
		}
	}
}

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
