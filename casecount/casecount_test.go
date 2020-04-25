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

var clientGetCallCounter int

func (m *mockClient) Get(url string) (*http.Response, error) {
	var csvStr string
	if clientGetCallCounter < 3 {
		csvStr = "Province/State,Country/Region,Lat,Long,1/22/20,1/23/20,1/24/20\n,Afghanistan,33.0,65.1,2,3,4\n,Albania,41.1533,20.1683,4,5,6\n,Algeria,28.0339,1.6596,7,8,9\n,US,37.0902,-95.7129,10,11,12"
	} else if clientGetCallCounter == 3 {
		csvStr = "UID,iso2,iso3,code3,FIPS,Admin2,Province_State,Country_Region,Lat,Long_,Combined_Key,1/22/20,1/23/20,1/24/20\n16,AS,ASM,16,60.0,,American Samoa,US,-14.270999999999999,-170.132,\"American Samoa, US\",4,5,6"
	} else {
		csvStr = "UID,iso2,iso3,code3,FIPS,Admin2,Province_State,Country_Region,Lat,Long_,Combined_Key,Population,1/22/20,1/23/20,1/24/20\n16,AS,ASM,16,60.0,,American Samoa,US,-14.270999999999999,-170.132,\"American Samoa, US\",55641,1,2,3"
	}
	clientGetCallCounter++
	r := ioutil.NopCloser(bytes.NewReader([]byte(csvStr)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

func getTestCacheData() map[string]map[string]CaseCounts {
	result := map[string]map[string]CaseCounts{
		"AF": map[string]CaseCounts{
			"": CaseCounts{
				Location{33.0, 65.1},
				[]CaseCount{
					CaseCount{"1/22/20", statistics{2, 2, 2}},
					CaseCount{"1/23/20", statistics{3, 3, 3}},
					CaseCount{"1/24/20", statistics{4, 4, 4}},
				},
			},
		},
		"AL": map[string]CaseCounts{
			"": CaseCounts{
				Location{41.1533, 20.1683},
				[]CaseCount{
					CaseCount{"1/22/20", statistics{4, 4, 4}},
					CaseCount{"1/23/20", statistics{5, 5, 5}},
					CaseCount{"1/24/20", statistics{6, 6, 6}},
				},
			},
		},
		"DZ": map[string]CaseCounts{
			"": CaseCounts{
				Location{28.0339, 1.6596},
				[]CaseCount{
					CaseCount{"1/22/20", statistics{7, 7, 7}},
					CaseCount{"1/23/20", statistics{8, 8, 8}},
					CaseCount{"1/24/20", statistics{9, 9, 9}},
				},
			},
		},
		"US": map[string]CaseCounts{
			"": CaseCounts{
				Location{37.0902, -95.7129},
				[]CaseCount{
					CaseCount{"1/22/20", statistics{0, 0, 10}},
					CaseCount{"1/23/20", statistics{0, 0, 11}},
					CaseCount{"1/24/20", statistics{0, 0, 12}},
				},
			},
			"American Samoa": CaseCounts{
				Location{-14.270999999999999, -170.132},
				[]CaseCount{
					CaseCount{"1/22/20", statistics{4, 1, 0}},
					CaseCount{"1/23/20", statistics{5, 2, 0}},
					CaseCount{"1/24/20", statistics{6, 3, 0}},
				},
			},
		},
	}
	return result
}

func TestUpdateCaseCounts(t *testing.T) {
	client = &mockClient{}
	clientGetCallCounter = 0
	UpdateCaseCounts()
	if firstDate.Format(dateformat.CasesDateFormat) != "1/22/20" {
		t.Errorf("Value of firstDate is incorrect, got: %s, want %s.", firstDate, "1/22/20")
	}
	if lastDate.Format(dateformat.CasesDateFormat) != "1/24/20" {
		t.Errorf("Value of lastDate is incorrect, got: %s, want %s.", lastDate, "1/24/20")
	}
	expectedCaseCounts := getTestCacheData()

	if len(caseCountsMap) != 4 {
		t.Errorf("Length of confirmedData is incorrect, got: %d, want %d.", len(caseCountsMap), 4)
	}
	verifyResultsCaseCountsMap(caseCountsMap, expectedCaseCounts, t)

	expectedAllAgg := map[string]map[string]CaseCountsAggregated{
		"AF": map[string]CaseCountsAggregated{
			"": CaseCountsAggregated{
				Location{33.0, 65.1},
				statistics{4, 4, 4},
			},
		},
		"AL": map[string]CaseCountsAggregated{
			"": CaseCountsAggregated{
				Location{41.1533, 20.1683},
				statistics{6, 6, 6},
			},
		},
		"DZ": map[string]CaseCountsAggregated{
			"": CaseCountsAggregated{
				Location{28.0339, 1.6596},
				statistics{9, 9, 9},
			},
		},
		"US": map[string]CaseCountsAggregated{
			"": CaseCountsAggregated{
				Location{37.0902, -95.7129},
				statistics{0, 0, 12},
			},
			"American Samoa": CaseCountsAggregated{
				Location{-14.270999999999999, -170.132},
				statistics{6, 3, 0},
			},
		},
	}
	caseCountsAgg, _ := GetCaseCounts("", "", "")
	verifyResultsCaseCountsAgg(caseCountsAgg, expectedAllAgg, t)

	expectedAllCountryAgg := map[string]CaseCountsAggregated{
		"AF": CaseCountsAggregated{
			Location{33.0, 65.1},
			statistics{4, 4, 4},
		},
		"AL": CaseCountsAggregated{
			Location{41.1533, 20.1683},
			statistics{6, 6, 6},
		},
		"DZ": CaseCountsAggregated{
			Location{28.0339, 1.6596},
			statistics{9, 9, 9},
		},
		"US": CaseCountsAggregated{
			Location{(-14.270999999999999 + 37.0902) / 2.0, (-95.7129 - 170.132) / 2},
			statistics{6, 3, 12},
		},
	}
	countryCaseCountsAgg, _ := GetCountryCaseCounts("", "", "")
	verifyResultsCountryCaseCountsAgg(countryCaseCountsAgg, expectedAllCountryAgg, t)

	caseCountsMap = nil
	stateAggregatedMap = nil
	countryAggregatedMap = nil
}

func TestGetCounts(t *testing.T) {
	client = &mockClient{}
	clientGetCallCounter = 0
	UpdateCaseCounts()
	expectedQueryAgg := map[string]map[string]CaseCountsAggregated{
		"AF": map[string]CaseCountsAggregated{
			"": CaseCountsAggregated{
				Location{33.0, 65.1},
				statistics{2, 2, 2},
			},
		},
		"AL": map[string]CaseCountsAggregated{
			"": CaseCountsAggregated{
				Location{41.1533, 20.1683},
				statistics{2, 2, 2},
			},
		},
		"DZ": map[string]CaseCountsAggregated{
			"": CaseCountsAggregated{
				Location{28.0339, 1.6596},
				statistics{2, 2, 2},
			},
		},
		"US": map[string]CaseCountsAggregated{
			"": CaseCountsAggregated{
				Location{37.0902, -95.7129},
				statistics{0, 0, 2},
			},
			"American Samoa": CaseCountsAggregated{
				Location{-14.270999999999999, -170.132},
				statistics{2, 2, 0},
			},
		},
	}
	caseCountsAgg, _ := GetCaseCounts("1/23/20", "1/24/20", "")
	verifyResultsCaseCountsAgg(caseCountsAgg, expectedQueryAgg, t)

	expectedQueryCountryAgg := map[string]CaseCountsAggregated{
		"AF": CaseCountsAggregated{
			Location{33.0, 65.1},
			statistics{3, 3, 3},
		},
		"AL": CaseCountsAggregated{
			Location{41.1533, 20.1683},
			statistics{5, 5, 5},
		},
		"DZ": CaseCountsAggregated{
			Location{28.0339, 1.6596},
			statistics{8, 8, 8},
		},
		"US": CaseCountsAggregated{
			Location{(-14.270999999999999 + 37.0902) / 2.0, (-95.7129 - 170.132) / 2.0},
			statistics{5, 2, 11},
		},
	}
	countryCaseCountsAgg, _ := GetCountryCaseCounts("1/22/20", "1/23/20", "")
	verifyResultsCountryCaseCountsAgg(countryCaseCountsAgg, expectedQueryCountryAgg, t)

	caseCountsMap = nil
	stateAggregatedMap = nil
	countryAggregatedMap = nil
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
	var input = []CaseCount{
		CaseCount{"a", statistics{2, 1, 0}},
		CaseCount{"b", statistics{4, 2, 1}},
		CaseCount{"c", statistics{7, 5, 3}},
	}

	tables := []struct {
		fromIndex         int
		toIndex           int
		expectedComfirmed int
		expectedDeaths    int
		expectedRecovered int
	}{
		{1, 2, 5, 4, 3},
		{-1, 2, 7, 5, 3},
		{1, 3, 5, 4, 3},
		{-2, -1, 0, 0, 0},
		{3, 4, 0, 0, 0},
		{2, 3, 3, 3, 2},
	}

	for _, table := range tables {
		confirmed, deaths, recovered := getStatisticsSum(input, table.fromIndex, table.toIndex)
		if confirmed != table.expectedComfirmed {
			t.Errorf("Confirmed was not correct, got: %d, want %d.", confirmed, table.expectedComfirmed)
		}
		if deaths != table.expectedDeaths {
			t.Errorf("Deaths was not 0, got: %d, want %d.", deaths, table.expectedDeaths)
		}
		if recovered != table.expectedRecovered {
			t.Errorf("Deaths was not 0, got: %d, want %d.", recovered, table.expectedRecovered)
		}
	}
}
