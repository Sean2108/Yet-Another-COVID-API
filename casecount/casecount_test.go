package casecount

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"yet-another-covid-map-api/dateformat"
)

type mockClient struct{}

func (m *mockClient) Get(url string) (*http.Response, error) {
	csvStr := "Province/State,Country/Region,Lat,Long,1/22/20,1/23/20,1/24/20\n,Afghanistan,33.0,65.1,2,3,4\n,Albania,41.1533,20.1683,4,5,6\n,Algeria,28.0339,1.6596,7,8,9"
	r := ioutil.NopCloser(bytes.NewReader([]byte(csvStr)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

func getTestCacheData() []CaseCounts {
	data1 := CaseCounts{stateInformation{"", "Afghanistan", 33.0, 65.1},
		[]caseCount{
			caseCount{"1/22/20", statistics{2, 2}},
			caseCount{"1/23/20", statistics{3, 3}},
			caseCount{"1/24/20", statistics{4, 4}},
		},
	}
	data2 := CaseCounts{stateInformation{"", "Albania", 41.1533, 20.1683},
		[]caseCount{
			caseCount{"1/22/20", statistics{4, 4}},
			caseCount{"1/23/20", statistics{5, 5}},
			caseCount{"1/24/20", statistics{6, 6}},
		},
	}
	data3 := CaseCounts{stateInformation{"", "Algeria", 28.0339, 1.6596},
		[]caseCount{
			caseCount{"1/22/20", statistics{7, 7}},
			caseCount{"1/23/20", statistics{8, 8}},
			caseCount{"1/24/20", statistics{9, 9}},
		},
	}
	return []CaseCounts{data1, data2, data3}
}

func TestUpdateCaseCounts(t *testing.T) {
	client = &mockClient{}
	UpdateCaseCounts()
	if firstDate.Format(dateformat.CasesDateFormat) != "1/22/20" {
		t.Errorf("Value of firstDate is incorrect, got: %s, want %s.", firstDate, "1/22/20")
	}
	if lastDate.Format(dateformat.CasesDateFormat) != "1/24/20" {
		t.Errorf("Value of lastDate is incorrect, got: %s, want %s.", lastDate, "1/24/20")
	}
	expectedCaseCounts := getTestCacheData()

	if len(caseCountsCache) != 3 {
		t.Errorf("Length of confirmedData is incorrect, got: %d, want %d.", len(caseCountsCache), 3)
	}
	verifyResultsCaseCountsArr(caseCountsCache, expectedCaseCounts, t)

	expectedAllAgg := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"", "Afghanistan", 33.0, 65.1}, statistics{4, 4}},
		CaseCountsAggregated{stateInformation{"", "Albania", 41.1533, 20.1683}, statistics{6, 6}},
		CaseCountsAggregated{stateInformation{"", "Algeria", 28.0339, 1.6596}, statistics{9, 9}},
	}
	caseCountsAgg, _ := GetCaseCounts("", "", "")
	verifyResultsCaseCountsAgg(caseCountsAgg, expectedAllAgg, t)

	expectedAllCountryAgg := []CountryCaseCountsAggregated{
		CountryCaseCountsAggregated{countryInformation{"Afghanistan", 33.0, 65.1}, statistics{4, 4}},
		CountryCaseCountsAggregated{countryInformation{"Albania", 41.1533, 20.1683}, statistics{6, 6}},
		CountryCaseCountsAggregated{countryInformation{"Algeria", 28.0339, 1.6596}, statistics{9, 9}},
	}
	countryCaseCountsAgg, _ := GetCountryCaseCounts("", "", "")
	verifyResultsCountryCaseCountsAgg(countryCaseCountsAgg, expectedAllCountryAgg, t)

	caseCountsCache = nil
	allAggregatedData = nil
	allCountriesAggregatedData = nil
}

func TestGetCounts(t *testing.T) {
	client = &mockClient{}
	UpdateCaseCounts()
	expectedAllAgg := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"", "Afghanistan", 33.0, 65.1}, statistics{2, 2}},
		CaseCountsAggregated{stateInformation{"", "Albania", 41.1533, 20.1683}, statistics{2, 2}},
		CaseCountsAggregated{stateInformation{"", "Algeria", 28.0339, 1.6596}, statistics{2, 2}},
	}
	caseCountsAgg, _ := GetCaseCounts("1/23/20", "1/24/20", "")
	verifyResultsCaseCountsAgg(caseCountsAgg, expectedAllAgg, t)

	expectedAllCountryAgg := []CountryCaseCountsAggregated{
		CountryCaseCountsAggregated{countryInformation{"Afghanistan", 33.0, 65.1}, statistics{3, 3}},
		CountryCaseCountsAggregated{countryInformation{"Albania", 41.1533, 20.1683}, statistics{5, 5}},
		CountryCaseCountsAggregated{countryInformation{"Algeria", 28.0339, 1.6596}, statistics{8, 8}},
	}
	countryCaseCountsAgg, _ := GetCountryCaseCounts("1/22/20", "1/23/20", "")
	verifyResultsCountryCaseCountsAgg(countryCaseCountsAgg, expectedAllCountryAgg, t)

	caseCountsCache = nil
	allAggregatedData = nil
	allCountriesAggregatedData = nil
}

func TestGetDaysBetweenDates(t *testing.T) {
	tables := []struct {
		startDate time.Time
		endDate   time.Time
		expected  int
	}{
		{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), 1},
		{time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC), time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), 1},
		{time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC), time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), 1},
		{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, 1, 11, 0, 0, 0, 0, time.UTC), 10},
	}

	for _, table := range tables {
		result := getDaysBetweenDates(table.startDate, table.endDate)
		if result != table.expected {
			t.Errorf("Difference between %s and %s was incorrect, got: %d, want: %d.", table.startDate, table.endDate, result, table.expected)
		}
	}
}

func TestGetStatisticsSum(t *testing.T) {
	var input = []caseCount{
		caseCount{"a", statistics{2, 1}},
		caseCount{"b", statistics{4, 2}},
		caseCount{"c", statistics{7, 5}},
	}

	tables := []struct {
		fromIndex         int
		toIndex           int
		expectedComfirmed int
		expectedDeaths    int
	}{
		{1, 2, 5, 4},
		{-1, 2, 7, 5},
		{1, 3, 5, 4},
		{-2, -1, 0, 0},
		{3, 4, 0, 0},
		{2, 3, 3, 3},
	}

	for _, table := range tables {
		confirmed, deaths := getStatisticsSum(input, table.fromIndex, table.toIndex)
		if confirmed != table.expectedComfirmed {
			t.Errorf("Confirmed was not correct, got: %d, want %d.", confirmed, table.expectedComfirmed)
		}
		if deaths != table.expectedDeaths {
			t.Errorf("Deaths was not 0, got: %d, want %d.", deaths, table.expectedDeaths)
		}
	}
}
