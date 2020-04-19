package casecount

import (
	"testing"
)

func (a *CaseCounts) equals(b CaseCounts) bool {
	if int(a.Lat) != int(b.Lat) || int(a.Long) != int(b.Long) {
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

func (a *CaseCountsAggregated) equals(b CaseCountsAggregated) bool {
	if int(a.Lat) != int(b.Lat) || int(a.Long) != int(b.Long) {
		return false
	}
	if a.statistics != b.statistics {
		return false
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

func verifyResultsCaseCountsMap(result map[string]map[string]CaseCounts, expectedData map[string]map[string]CaseCounts, t *testing.T) {
	for country, countryInfo := range result {
		for state, stateInfo := range countryInfo {
			if !stateInfo.equals(expectedData[country][state]) {
				t.Errorf("Result data is incorrect, got: %+v, want %+v.", stateInfo, expectedData[country][state])
			}
		}
	}
}

func verifyResultsCountryCaseCountsMap(result map[string]CaseCounts, expectedData map[string]CaseCounts, t *testing.T) {
	for country, countryInfo := range result {
		if !countryInfo.equals(expectedData[country]) {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", countryInfo, expectedData[country])
		}
	}
}

func verifyResultsCaseCountsAgg(result map[string]map[string]CaseCountsAggregated, expectedData map[string]map[string]CaseCountsAggregated, t *testing.T) {
	for country, countryInfo := range result {
		for state, stateInfo := range countryInfo {
			if !stateInfo.equals(expectedData[country][state]) {
				t.Errorf("Result data is incorrect, got: %+v, want %+v.", stateInfo, expectedData[country][state])
			}
		}
	}
}

func verifyResultsCountryCaseCountsAgg(result map[string]CaseCountsAggregated, expectedData map[string]CaseCountsAggregated, t *testing.T) {
	for country, countryInfo := range result {
		if !countryInfo.equals(expectedData[country]) {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", countryInfo, expectedData[country])
		}
	}
}
