package casecount

import (
	"strings"
	"testing"
)

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

func TestAggregateDataBetweenDates_QueryDatesBeforeValidRange(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/20/20", "1/21/20", "")
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"Beijing", "China", 40.1824, 116.4142}, statistics{0, 0}},
		CaseCountsAggregated{stateInformation{"Hubei", "China", 30.9756, 112.2707}, statistics{0, 0}},
		CaseCountsAggregated{stateInformation{"Shanghai", "China", 31.202, 121.4491}, statistics{0, 0}},
		CaseCountsAggregated{stateInformation{"", "Singapore", 1.2833, 103.8333}, statistics{0, 0}},
		CaseCountsAggregated{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004}, statistics{0, 0}},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryDatesAfterValidRange(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/28/20", "1/29/20", "")
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := []CaseCountsAggregated{
		CaseCountsAggregated{stateInformation{"Beijing", "China", 40.1824, 116.4142}, statistics{0, 0}},
		CaseCountsAggregated{stateInformation{"Hubei", "China", 30.9756, 112.2707}, statistics{0, 0}},
		CaseCountsAggregated{stateInformation{"Shanghai", "China", 31.202, 121.4491}, statistics{0, 0}},
		CaseCountsAggregated{stateInformation{"", "Singapore", 1.2833, 103.8333}, statistics{0, 0}},
		CaseCountsAggregated{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004}, statistics{0, 0}},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryDatesBeforeAndAfter_ShouldReturnAll(t *testing.T) {
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
	result, err := aggregateDataBetweenDates("", "", "Siingapore")
	if len(result) != 0 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 0)
	}
	if !strings.Contains(err.Error(), "Singapore") {
		t.Errorf("Error message is incorrect, got: %s, want %s.", err.Error(), "string containing Singapore")
	}
	result, err = aggregateDataBetweenDates("", "", "chain")
	if len(result) != 0 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 0)
	}
	if !strings.Contains(err.Error(), "China") {
		t.Errorf("Error message is incorrect, got: %s, want %s.", err.Error(), "string containing China")
	}
	result, err = aggregateDataBetweenDates("", "", "UnitedKingdom")
	if len(result) != 0 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 0)
	}
	if !strings.Contains(err.Error(), "United Kingdom") {
		t.Errorf("Error message is incorrect, got: %s, want %s.", err.Error(), "string containing United Kingdom")
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
