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

func verifyResultsCaseCountArr(expectedData []CaseCount, result []CaseCount, t *testing.T) {
	for i, item := range result {
		if item != expectedData[i] {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", item, expectedData[i])
		}
	}
}

func verifyResultsCaseCountsMap(expectedData map[string]CountryWithStates, result map[string]CountryWithStates, t *testing.T) {
	for country, countryInfo := range result {
		if countryInfo.Name != expectedData[country].Name {
			t.Errorf("Country names are incorrect, got: %+v, want %+v.", countryInfo.Name, expectedData[country].Name)
		}
		for state, stateInfo := range countryInfo.States {
			if !stateInfo.equals(expectedData[country].States[state]) {
				t.Errorf("Result data is incorrect, got: %+v, want %+v.", stateInfo, expectedData[country].States[state])
			}
		}
	}
}

func verifyResultsCountryCaseCountsMap(expectedData map[string]Country, result map[string]Country, t *testing.T) {
	for country, countryInfo := range result {
		if countryInfo.Name != expectedData[country].Name {
			t.Errorf("Country names are incorrect, got: %+v, want %+v.", countryInfo.Name, expectedData[country].Name)
		}
		if !countryInfo.CaseCounts.equals(expectedData[country].CaseCounts) {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", countryInfo, expectedData[country])
		}
	}
}

func verifyResultsCaseCountsAgg(expectedData map[string]CountryWithStatesAggregated, result map[string]CountryWithStatesAggregated, t *testing.T) {
	for country, countryInfo := range result {
		if countryInfo.Name != expectedData[country].Name {
			t.Errorf("Country names are incorrect, got: %+v, want %+v.", countryInfo.Name, expectedData[country].Name)
		}
		for state, stateInfo := range countryInfo.States {
			if !stateInfo.equals(expectedData[country].States[state]) {
				t.Errorf("Result data is incorrect, got: %+v, want %+v.", stateInfo, expectedData[country].States[state])
			}
		}
	}
}

func verifyResultsCountryCaseCountsAgg(expectedData map[string]CountryAggregated, result map[string]CountryAggregated, t *testing.T) {
	for country, countryInfo := range result {
		if countryInfo.Name != expectedData[country].Name {
			t.Errorf("Country names are incorrect, got: %+v, want %+v.", countryInfo.Name, expectedData[country].Name)
		}
		if !countryInfo.CaseCountsAggregated.equals(expectedData[country].CaseCountsAggregated) {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", countryInfo, expectedData[country])
		}
	}
}
