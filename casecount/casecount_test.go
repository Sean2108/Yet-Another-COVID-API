package casecount

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"yet-another-covid-map-api/dateformat"
	"yet-another-covid-map-api/utils"
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

func getTestCacheData() map[string]CountryWithStates {
	result := map[string]CountryWithStates{
		"AF": CountryWithStates{
			Name: "Afghanistan",
			States: map[string]CaseCounts{
				"": CaseCounts{
					LocationAndPopulation{33.0, 65.1, 5000},
					[]CaseCount{
						CaseCount{"1/22/20", statistics{2, 2, 2}},
						CaseCount{"1/23/20", statistics{3, 3, 3}},
						CaseCount{"1/24/20", statistics{4, 4, 4}},
					},
				},
			},
		},
		"AL": CountryWithStates{
			Name: "Albania",
			States: map[string]CaseCounts{
				"": CaseCounts{
					LocationAndPopulation{41.1533, 20.1683, 3000},
					[]CaseCount{
						CaseCount{"1/22/20", statistics{4, 4, 4}},
						CaseCount{"1/23/20", statistics{5, 5, 5}},
						CaseCount{"1/24/20", statistics{6, 6, 6}},
					},
				},
			},
		},
		"DZ": CountryWithStates{
			Name: "Algeria",
			States: map[string]CaseCounts{
				"": CaseCounts{
					LocationAndPopulation{28.0339, 1.6596, 6000},
					[]CaseCount{
						CaseCount{"1/22/20", statistics{7, 7, 7}},
						CaseCount{"1/23/20", statistics{8, 8, 8}},
						CaseCount{"1/24/20", statistics{9, 9, 9}},
					},
				},
			},
		},
		"US": CountryWithStates{
			Name: "US",
			States: map[string]CaseCounts{
				"": CaseCounts{
					LocationAndPopulation{37.0902, -95.7129, 300000},
					[]CaseCount{
						CaseCount{"1/22/20", statistics{0, 0, 10}},
						CaseCount{"1/23/20", statistics{0, 0, 11}},
						CaseCount{"1/24/20", statistics{0, 0, 12}},
					},
				},
				"American Samoa": CaseCounts{
					LocationAndPopulation{-14.270999999999999, -170.132, 40000},
					[]CaseCount{
						CaseCount{"1/22/20", statistics{4, 1, 0}},
						CaseCount{"1/23/20", statistics{5, 2, 0}},
						CaseCount{"1/24/20", statistics{6, 3, 0}},
					},
				},
			},
		},
	}
	return result
}

func setupTest() {
	utils.AbbreviationToCountry = map[string]string{
		"AF": "Afghanistan",
		"AL": "Albania",
		"DZ": "Algeria",
		"US": "US",
		"CN": "China",
		"SG": "Singapore",
		"GB": "United Kingdom",
	}
	utils.CountryToAbbreviation = map[string]string{
		"Afghanistan":    "AF",
		"Albania":        "AL",
		"Algeria":        "DZ",
		"US":             "US",
		"China":          "CN",
		"Singapore":      "SG",
		"United Kingdom": "GB",
	}
	utils.StatePopulationLookup = map[string]map[string]int{
		"AF": map[string]int{
			"": 5000,
		},
		"AL": map[string]int{
			"": 3000,
		},
		"DZ": map[string]int{
			"": 6000,
		},
		"US": map[string]int{
			"":               300000,
			"American Samoa": 40000,
		},
		"CN": map[string]int{
			"Hubei":    300000,
			"Shanghai": 40000,
			"Beijing":  50000,
		},
		"SG": map[string]int{
			"": 6000,
		},
		"GB": map[string]int{
			"London": 7000,
		},
	}
}

func init() {
	setupTest()
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

	expectedAllAgg := map[string]CountryWithStatesAggregated{
		"AF": CountryWithStatesAggregated{
			Name: "Afghanistan",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{33.0, 65.1, 5000},
					statistics{4, 4, 4},
				},
			},
		},
		"AL": CountryWithStatesAggregated{
			Name: "Albania",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{41.1533, 20.1683, 3000},
					statistics{6, 6, 6},
				},
			},
		},
		"DZ": CountryWithStatesAggregated{
			Name: "Algeria",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{28.0339, 1.6596, 6000},
					statistics{9, 9, 9},
				},
			},
		},
		"US": CountryWithStatesAggregated{
			Name: "US",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{37.0902, -95.7129, 300000},
					statistics{0, 0, 12},
				},
				"American Samoa": CaseCountsAggregated{
					LocationAndPopulation{-14.270999999999999, -170.132, 40000},
					statistics{6, 3, 0},
				},
			},
		},
	}
	caseCountsAgg, _ := GetCaseCounts("", "", "")
	verifyResultsCaseCountsAgg(caseCountsAgg, expectedAllAgg, t)

	expectedAllCountryAgg := map[string]CountryAggregated{
		"AF": CountryAggregated{
			"Afghanistan",
			CaseCountsAggregated{
				LocationAndPopulation{33.0, 65.1, 5000},
				statistics{4, 4, 4},
			},
		},
		"AL": CountryAggregated{
			"Albania",
			CaseCountsAggregated{
				LocationAndPopulation{41.1533, 20.1683, 3000},
				statistics{6, 6, 6},
			},
		},
		"DZ": CountryAggregated{
			"Algeria",
			CaseCountsAggregated{
				LocationAndPopulation{28.0339, 1.6596, 6000},
				statistics{9, 9, 9},
			},
		},
		"US": CountryAggregated{
			"US",
			CaseCountsAggregated{
				LocationAndPopulation{37.0902, -95.7129, 340000},
				statistics{6, 3, 12},
			},
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

	expectedQueryAgg := map[string]CountryWithStatesAggregated{
		"AF": CountryWithStatesAggregated{
			Name: "Afghanistan",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{33.0, 65.1, 5000},
					statistics{2, 2, 2},
				},
			},
		},
		"AL": CountryWithStatesAggregated{
			Name: "Albania",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{41.1533, 20.1683, 3000},
					statistics{2, 2, 2},
				},
			},
		},
		"DZ": CountryWithStatesAggregated{
			Name: "Algeria",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{28.0339, 1.6596, 6000},
					statistics{2, 2, 2},
				},
			},
		},
		"US": CountryWithStatesAggregated{
			Name: "US",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{37.0902, -95.7129, 300000},
					statistics{0, 0, 2},
				},
				"American Samoa": CaseCountsAggregated{
					LocationAndPopulation{-14.270999999999999, -170.132, 40000},
					statistics{2, 2, 0},
				},
			},
		},
	}
	caseCountsAgg, _ := GetCaseCounts("1/23/20", "1/24/20", "")
	verifyResultsCaseCountsAgg(caseCountsAgg, expectedQueryAgg, t)

	expectedQueryCountryAgg := map[string]CountryAggregated{
		"AF": CountryAggregated{
			"Afghanistan",
			CaseCountsAggregated{
				LocationAndPopulation{33.0, 65.1, 5000},
				statistics{3, 3, 3},
			},
		},
		"AL": CountryAggregated{
			"Albania",
			CaseCountsAggregated{
				LocationAndPopulation{41.1533, 20.1683, 3000},
				statistics{5, 5, 5},
			},
		},
		"DZ": CountryAggregated{
			"Algeria",
			CaseCountsAggregated{
				LocationAndPopulation{28.0339, 1.6596, 6000},
				statistics{8, 8, 8},
			},
		},
		"US": CountryAggregated{
			"US",
			CaseCountsAggregated{
				LocationAndPopulation{37.0902, -95.7129, 340000},
				statistics{5, 2, 11},
			},
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
