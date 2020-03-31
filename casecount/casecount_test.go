package casecount

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"
)

type ByCountryAndStateForCaseCounts []caseCounts

func (a ByCountryAndStateForCaseCounts) Len() int {
	return len(a)
}

func (a ByCountryAndStateForCaseCounts) Less(i, j int) bool {
	if a[i].Country == a[j].Country {
		return a[i].State < a[j].State
	}
	return a[i].Country < a[j].Country
}

func (a ByCountryAndStateForCaseCounts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ByCountryAndStateAgg []CaseCountsAggregated

func (a ByCountryAndStateAgg) Len() int {
	return len(a)
}

func (a ByCountryAndStateAgg) Less(i, j int) bool {
	if a[i].Country == a[j].Country {
		return a[i].State < a[j].State
	}
	return a[i].Country < a[j].Country
}

func (a ByCountryAndStateAgg) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ByCountryAgg []CountryCaseCountsAggregated

func (a ByCountryAgg) Len() int {
	return len(a)
}

func (a ByCountryAgg) Less(i, j int) bool {
	return a[i].Country < a[j].Country
}

func (a ByCountryAgg) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type mockClient struct{}

var mockedGet func(url string) (*http.Response, error)

func (m *mockClient) Get(url string) (*http.Response, error) {
	return mockedGet(url)
}

func TestUpdateCaseCounts(t *testing.T) {
	client = &mockClient{}
	mockedGet = func(url string) (*http.Response, error) {
		csvStr := "Province/State,Country/Region,Lat,Long,1/22/20,1/23/20,1/24/20\n,Afghanistan,33.0,65.1,2,3,4\n,Albania,41.1533,20.1683,4,5,6\n,Algeria,28.0339,1.6596,7,8,9"
		r := ioutil.NopCloser(bytes.NewReader([]byte(csvStr)))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	UpdateCaseCounts()
	if firstDate.Format(inputDateFormat) != "1/22/20" {
		t.Errorf("Value of firstDate is incorrect, got: %s, want %s.", firstDate, "1/22/20")
	}
	if lastDate.Format(inputDateFormat) != "1/24/20" {
		t.Errorf("Value of lastDate is incorrect, got: %s, want %s.", lastDate, "1/24/20")
	}
	expectedData1 := caseCounts{stateInformation{"", "Afghanistan", 33.0, 65.1},
		[]caseCount{
			caseCount{"1/22/20", statistics{2, 2}},
			caseCount{"1/23/20", statistics{3, 3}},
			caseCount{"1/24/20", statistics{4, 4}},
		},
	}
	expectedData2 := caseCounts{stateInformation{"", "Albania", 41.1533, 20.1683},
		[]caseCount{
			caseCount{"1/22/20", statistics{4, 4}},
			caseCount{"1/23/20", statistics{5, 5}},
			caseCount{"1/24/20", statistics{6, 6}},
		},
	}
	expectedData3 := caseCounts{stateInformation{"", "Algeria", 28.0339, 1.6596},
		[]caseCount{
			caseCount{"1/22/20", statistics{7, 7}},
			caseCount{"1/23/20", statistics{8, 8}},
			caseCount{"1/24/20", statistics{9, 9}},
		},
	}
	expectedCaseCounts := []caseCounts{expectedData1, expectedData2, expectedData3}

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
	confirmed, deaths := getStatisticsSum(input, 1, 2)
	if confirmed != 5 {
		t.Errorf("Confirmed was not correct, got: %d, want %d.", confirmed, 5)
	}
	if deaths != 4 {
		t.Errorf("Deaths was not 0, got: %d, want %d.", deaths, 4)
	}
}

func (a *caseCounts) Equals(b caseCounts) bool {
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

func getTestCaseCounts() []caseCounts {
	expectedBeijingData := caseCounts{stateInformation{"Beijing", "China", 40.1824, 116.4142},
		[]caseCount{
			caseCount{"1/22/20", statistics{50, 10}},
			caseCount{"1/23/20", statistics{200, 87}},
			caseCount{"1/24/20", statistics{800, 125}},
			caseCount{"1/25/20", statistics{1020, 142}},
			caseCount{"1/26/20", statistics{1110, 145}},
			caseCount{"1/27/20", statistics{1235, 152}},
		},
	}
	expectedHubeiData := caseCounts{stateInformation{"Hubei", "China", 30.9756, 112.2707},
		[]caseCount{
			caseCount{"1/22/20", statistics{100, 20}},
			caseCount{"1/23/20", statistics{1000, 100}},
			caseCount{"1/24/20", statistics{1800, 105}},
			caseCount{"1/25/20", statistics{2020, 150}},
			caseCount{"1/26/20", statistics{2110, 175}},
			caseCount{"1/27/20", statistics{2111, 230}},
		},
	}
	expectedShanghaiData := caseCounts{stateInformation{"Shanghai", "China", 31.202, 121.4491},
		[]caseCount{
			caseCount{"1/22/20", statistics{10, 5}},
			caseCount{"1/23/20", statistics{45, 8}},
			caseCount{"1/24/20", statistics{89, 20}},
			caseCount{"1/25/20", statistics{126, 25}},
			caseCount{"1/26/20", statistics{400, 42}},
			caseCount{"1/27/20", statistics{532, 55}},
		},
	}
	expectedSingaporeData := caseCounts{stateInformation{"", "Singapore", 1.2833, 103.8333},
		[]caseCount{
			caseCount{"1/22/20", statistics{1, 0}},
			caseCount{"1/23/20", statistics{3, 2}},
			caseCount{"1/24/20", statistics{6, 4}},
			caseCount{"1/25/20", statistics{10, 5}},
			caseCount{"1/26/20", statistics{15, 8}},
			caseCount{"1/27/20", statistics{23, 10}},
		},
	}
	expectedLondonData := caseCounts{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004},
		[]caseCount{
			caseCount{"1/22/20", statistics{1, 0}},
			caseCount{"1/23/20", statistics{6, 1}},
			caseCount{"1/24/20", statistics{8, 3}},
			caseCount{"1/25/20", statistics{9, 6}},
			caseCount{"1/26/20", statistics{20, 6}},
			caseCount{"1/27/20", statistics{28, 9}},
		},
	}
	return []caseCounts{expectedBeijingData, expectedHubeiData, expectedShanghaiData, expectedSingaporeData, expectedLondonData}
}

func verifyResultsCaseCountsArr(result []caseCounts, expectedData []caseCounts, t *testing.T) {
	sort.Sort(ByCountryAndStateForCaseCounts(result))
	for i, item := range result {
		if !item.Equals(expectedData[i]) {
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
		if item != expectedData[i] {
			t.Errorf("Result data is incorrect, got: %+v, want %+v.", item, expectedData[i])
		}
	}
}

func TestExtractCaseCounts(t *testing.T) {
	confirmedData := [][]string{
		{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"},
		{"Beijing", "China", "40.1824", "116.4142", "50", "200", "800", "1020", "1110", "1235"},
		{"Hubei", "China", "30.9756", "112.2707", "100", "1000", "1800", "2020", "2110", "2111"},
		{"Shanghai", "China", "31.202", "121.4491", "10", "45", "89", "126", "400", "532"},
		{"", "Singapore", "1.2833", "103.8333", "1", "3", "6", "10", "15", "23"},
		{"London", "United Kingdom", "55.3781", "-3.4360000000000004", "1", "6", "8", "9", "20", "28"},
	}
	deathsData := [][]string{
		{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"},
		{"Beijing", "China", "40.1824", "116.4142", "10", "87", "125", "142", "145", "152"},
		{"Hubei", "China", "30.9756", "112.2707", "20", "100", "105", "150", "175", "230"},
		{"Shanghai", "China", "31.202", "121.4491", "5", "8", "20", "25", "42", "55"},
		{"", "Singapore", "1.2833", "103.8333", "0", "2", "4", "5", "8", "10"},
		{"London", "United Kingdom", "55.3781", "-3.4360000000000004", "0", "1", "3", "6", "6", "9"},
	}
	headerRow := confirmedData[0]
	extractCaseCounts(headerRow, confirmedData, deathsData)
	result := caseCountsCache
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := getTestCaseCounts()
	verifyResultsCaseCountsArr(result, expectedData, t)
}

func TestAggregateDataBetweenDates_AllDates(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result := allAggregatedData
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"Beijing", "China", 40.1824, 116.4142}, statistics{1235, 152}},
		CaseCountsAggregated{stateInformation{"Hubei", "China", 30.9756, 112.2707}, statistics{2111, 230}},
		CaseCountsAggregated{stateInformation{"Shanghai", "China", 31.202, 121.4491}, statistics{532, 55}},
		CaseCountsAggregated{stateInformation{"", "Singapore", 1.2833, 103.8333}, statistics{23, 10}},
		CaseCountsAggregated{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004}, statistics{28, 9}},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryDates(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/24/20", "1/26/20", "")
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"Beijing", "China", 40.1824, 116.4142}, statistics{910, 58}},
		CaseCountsAggregated{stateInformation{"Hubei", "China", 30.9756, 112.2707}, statistics{1110, 75}},
		CaseCountsAggregated{stateInformation{"Shanghai", "China", 31.202, 121.4491}, statistics{355, 34}},
		CaseCountsAggregated{stateInformation{"", "Singapore", 1.2833, 103.8333}, statistics{12, 6}},
		CaseCountsAggregated{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004}, statistics{14, 5}},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryCountry(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("", "", "Singapore")
	if len(result) != 1 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 1)
	}
	expectedData := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"", "Singapore", 1.2833, 103.8333}, statistics{23, 10}},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryCountry_wEiRdCaSiNg(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("", "", "sInGaPoRe")
	if len(result) != 1 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 1)
	}
	expectedData := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"", "Singapore", 1.2833, 103.8333}, statistics{23, 10}},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryCountry_SpellingMistakeInCountryName(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, err := aggregateDataBetweenDates("", "", "Sngapore")
	if len(result) != 0 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 0)
	}
	if !strings.Contains(err.Error(), "Singapore") {
		t.Errorf("Error message is incorrect, got: %s, want %s.", err.Error(), "string containing Singapore")
	}
}

func TestAggregateDataBetweenDates_QueryDates_FromIsOutOfRange(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/21/20", "1/26/20", "")
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"Beijing", "China", 40.1824, 116.4142}, statistics{1110, 145}},
		CaseCountsAggregated{stateInformation{"Hubei", "China", 30.9756, 112.2707}, statistics{2110, 175}},
		CaseCountsAggregated{stateInformation{"Shanghai", "China", 31.202, 121.4491}, statistics{400, 42}},
		CaseCountsAggregated{stateInformation{"", "Singapore", 1.2833, 103.8333}, statistics{15, 8}},
		CaseCountsAggregated{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004}, statistics{20, 6}},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryDates_ToIsOutOfRange(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/24/20", "1/28/20", "")
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"Beijing", "China", 40.1824, 116.4142}, statistics{1035, 65}},
		CaseCountsAggregated{stateInformation{"Hubei", "China", 30.9756, 112.2707}, statistics{1111, 130}},
		CaseCountsAggregated{stateInformation{"Shanghai", "China", 31.202, 121.4491}, statistics{487, 47}},
		CaseCountsAggregated{stateInformation{"", "Singapore", 1.2833, 103.8333}, statistics{20, 8}},
		CaseCountsAggregated{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004}, statistics{22, 8}},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryDates_FromAndToBothOutOfRange(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/21/20", "1/28/20", "")
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"Beijing", "China", 40.1824, 116.4142}, statistics{1235, 152}},
		CaseCountsAggregated{stateInformation{"Hubei", "China", 30.9756, 112.2707}, statistics{2111, 230}},
		CaseCountsAggregated{stateInformation{"Shanghai", "China", 31.202, 121.4491}, statistics{532, 55}},
		CaseCountsAggregated{stateInformation{"", "Singapore", 1.2833, 103.8333}, statistics{23, 10}},
		CaseCountsAggregated{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004}, statistics{28, 9}},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateCountryDataFromStatesAggregate_AllDates(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result := allCountriesAggregatedData
	if len(result) != 3 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 3)
	}
	expectedData := []CountryCaseCountsAggregated{
		CountryCaseCountsAggregated{countryInformation{"China", (40.1824 + 30.9756 + 31.202) / 3.0, (116.4142 + 112.2707 + 121.4491) / 3.0}, statistics{3878, 437}},
		CountryCaseCountsAggregated{countryInformation{"Singapore", 1.2833, 103.8333}, statistics{23, 10}},
		CountryCaseCountsAggregated{countryInformation{"United Kingdom", 55.3781, -3.4360000000000004}, statistics{28, 9}},
	}
	verifyResultsCountryCaseCountsAgg(result, expectedData, t)
}

func TestAggregateCountryDataFromStatesAggregate_QueryDates(t *testing.T) {
	input := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"Beijing", "China", 40.1824, 116.4142}, statistics{910, 58}},
		CaseCountsAggregated{stateInformation{"Hubei", "China", 30.9756, 112.2707}, statistics{1110, 75}},
		CaseCountsAggregated{stateInformation{"Shanghai", "China", 31.202, 121.4491}, statistics{355, 34}},
		CaseCountsAggregated{stateInformation{"", "Singapore", 1.2833, 103.8333}, statistics{12, 6}},
		CaseCountsAggregated{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004}, statistics{14, 5}},
	}
	result := aggregateCountryDataFromStatesAggregate(input)
	if len(result) != 3 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 3)
	}
	expectedData := []CountryCaseCountsAggregated{
		CountryCaseCountsAggregated{countryInformation{"China", (40.1824 + 30.9756 + 31.202) / 3.0, (116.4142 + 112.2707 + 121.4491) / 3.0}, statistics{2375, 167}},
		CountryCaseCountsAggregated{countryInformation{"Singapore", 1.2833, 103.8333}, statistics{12, 6}},
		CountryCaseCountsAggregated{countryInformation{"United Kingdom", 55.3781, -3.4360000000000004}, statistics{14, 5}},
	}
	verifyResultsCountryCaseCountsAgg(result, expectedData, t)
}
