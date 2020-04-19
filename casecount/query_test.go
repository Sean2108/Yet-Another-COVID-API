package casecount

import (
	"strings"
	"testing"
)

func getTestCaseCounts() []CaseCounts {
	expectedBeijingData := CaseCounts{stateInformation{"Beijing", "China", 40.1824, 116.4142},
		[]CaseCount{
			CaseCount{"1/22/20", statistics{50, 10}},
			CaseCount{"1/23/20", statistics{200, 87}},
			CaseCount{"1/24/20", statistics{800, 125}},
			CaseCount{"1/25/20", statistics{1020, 142}},
			CaseCount{"1/26/20", statistics{1110, 145}},
			CaseCount{"1/27/20", statistics{1235, 152}},
		},
	}
	expectedHubeiData := CaseCounts{stateInformation{"Hubei", "China", 30.9756, 112.2707},
		[]CaseCount{
			CaseCount{"1/22/20", statistics{100, 20}},
			CaseCount{"1/23/20", statistics{1000, 100}},
			CaseCount{"1/24/20", statistics{1800, 105}},
			CaseCount{"1/25/20", statistics{2020, 150}},
			CaseCount{"1/26/20", statistics{2110, 175}},
			CaseCount{"1/27/20", statistics{2111, 230}},
		},
	}
	expectedShanghaiData := CaseCounts{stateInformation{"Shanghai", "China", 31.202, 121.4491},
		[]CaseCount{
			CaseCount{"1/22/20", statistics{10, 5}},
			CaseCount{"1/23/20", statistics{45, 8}},
			CaseCount{"1/24/20", statistics{89, 20}},
			CaseCount{"1/25/20", statistics{126, 25}},
			CaseCount{"1/26/20", statistics{400, 42}},
			CaseCount{"1/27/20", statistics{532, 55}},
		},
	}
	expectedSingaporeData := CaseCounts{stateInformation{"", "Singapore", 1.2833, 103.8333},
		[]CaseCount{
			CaseCount{"1/22/20", statistics{1, 0}},
			CaseCount{"1/23/20", statistics{3, 2}},
			CaseCount{"1/24/20", statistics{6, 4}},
			CaseCount{"1/25/20", statistics{10, 5}},
			CaseCount{"1/26/20", statistics{15, 8}},
			CaseCount{"1/27/20", statistics{23, 10}},
		},
	}
	expectedLondonData := CaseCounts{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004},
		[]CaseCount{
			CaseCount{"1/22/20", statistics{1, 0}},
			CaseCount{"1/23/20", statistics{6, 1}},
			CaseCount{"1/24/20", statistics{8, 3}},
			CaseCount{"1/25/20", statistics{9, 6}},
			CaseCount{"1/26/20", statistics{20, 6}},
			CaseCount{"1/27/20", statistics{28, 9}},
		},
	}
	return []CaseCounts{expectedBeijingData, expectedHubeiData, expectedShanghaiData, expectedSingaporeData, expectedLondonData}
}

func getTestCaseCountsWithoutFirstAndLastDay() []CaseCounts {
	expectedBeijingData := CaseCounts{stateInformation{"Beijing", "China", 40.1824, 116.4142},
		[]CaseCount{
			CaseCount{"1/23/20", statistics{200, 87}},
			CaseCount{"1/24/20", statistics{800, 125}},
			CaseCount{"1/25/20", statistics{1020, 142}},
			CaseCount{"1/26/20", statistics{1110, 145}},
		},
	}
	expectedHubeiData := CaseCounts{stateInformation{"Hubei", "China", 30.9756, 112.2707},
		[]CaseCount{
			CaseCount{"1/23/20", statistics{1000, 100}},
			CaseCount{"1/24/20", statistics{1800, 105}},
			CaseCount{"1/25/20", statistics{2020, 150}},
			CaseCount{"1/26/20", statistics{2110, 175}},
		},
	}
	expectedShanghaiData := CaseCounts{stateInformation{"Shanghai", "China", 31.202, 121.4491},
		[]CaseCount{
			CaseCount{"1/23/20", statistics{45, 8}},
			CaseCount{"1/24/20", statistics{89, 20}},
			CaseCount{"1/25/20", statistics{126, 25}},
			CaseCount{"1/26/20", statistics{400, 42}},
		},
	}
	expectedSingaporeData := CaseCounts{stateInformation{"", "Singapore", 1.2833, 103.8333},
		[]CaseCount{
			CaseCount{"1/23/20", statistics{3, 2}},
			CaseCount{"1/24/20", statistics{6, 4}},
			CaseCount{"1/25/20", statistics{10, 5}},
			CaseCount{"1/26/20", statistics{15, 8}},
		},
	}
	expectedLondonData := CaseCounts{stateInformation{"London", "United Kingdom", 55.3781, -3.4360000000000004},
		[]CaseCount{
			CaseCount{"1/23/20", statistics{6, 1}},
			CaseCount{"1/24/20", statistics{8, 3}},
			CaseCount{"1/25/20", statistics{9, 6}},
			CaseCount{"1/26/20", statistics{20, 6}},
		},
	}
	return []CaseCounts{expectedBeijingData, expectedHubeiData, expectedShanghaiData, expectedSingaporeData, expectedLondonData}
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
	if len(result) != 0 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 0)
	}
}

func TestAggregateDataBetweenDates_QueryDatesAfterValidRange(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/28/20", "1/29/20", "")
	if len(result) != 0 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 0)
	}
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

func TestAggregateDataBetweenDates_QueryFromDateAfterToDate(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	_, err := aggregateDataBetweenDates("1/24/20", "1/23/20", "China")
	if err == nil {
		t.Error("Error message should be returned.")
	}
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

func TestAggregateDataPerDay_AllDates(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCaseCountsWithDayData("", "", "")
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := caseCountsCache
	verifyResultsCaseCountsArr(result, expectedData, t)
}

func TestAggregateDataPerDay_QueryDates(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCaseCountsWithDayData("1/23/20", "1/26/20", "")
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := getTestCaseCountsWithoutFirstAndLastDay()
	verifyResultsCaseCountsArr(result, expectedData, t)
}

func TestAggregateDataPerDay_BeforeAndAfterShouldReturnAll(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCaseCountsWithDayData("1/21/20", "1/28/20", "")
	if len(result) != 5 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 5)
	}
	expectedData := caseCountsCache
	verifyResultsCaseCountsArr(result, expectedData, t)
}

func TestAggregateDataPerDay_CountryQuery(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCaseCountsWithDayData("", "", "China")
	if len(result) != 3 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 3)
	}
	expectedData := getTestCaseCounts()[0:3]
	verifyResultsCaseCountsArr(result, expectedData, t)
}

func TestAggregateDataPerDay_CountryQueryWithSpellingError(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, err := GetCaseCountsWithDayData("", "", "Chiina")
	if len(result) != 0 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 0)
	}
	if !strings.Contains(err.Error(), "China") {
		t.Errorf("Error message is incorrect, got: %s, want %s.", err.Error(), "string containing China")
	}
}

func TestAggregateDataPerDay_QueryFromDateAfterToDate(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	_, err := GetCaseCountsWithDayData("1/24/20", "1/23/20", "China")
	if err == nil {
		t.Error("Error message should be returned.")
	}
}

func TestCountryAggregateDataPerDay_AllDates(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCountryCaseCountsWithDayData("", "", "")
	if len(result) != 3 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 3)
	}
	expectedData := countryCaseCountsCache
	verifyResultsCountryCaseCountsArr(result, expectedData, t)
}

func TestCountryAggregateDataPerDay_QueryDates(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCountryCaseCountsWithDayData("1/23/20", "1/26/20", "")
	if len(result) != 3 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 3)
	}
	expectedData := []CountryCaseCounts{
		CountryCaseCounts{countryInformation{"China", (40.1824 + 30.9756 + 31.202) / 3.0, (116.4142 + 112.2707 + 121.4491) / 3.0},
			[]CaseCount{
				CaseCount{"1/23/20", statistics{1245, 195}},
				CaseCount{"1/24/20", statistics{2689, 250}},
				CaseCount{"1/25/20", statistics{3166, 317}},
				CaseCount{"1/26/20", statistics{3620, 362}},
			}},
		CountryCaseCounts{countryInformation{"Singapore", 1.2833, 103.8333}, []CaseCount{
			CaseCount{"1/23/20", statistics{3, 2}},
			CaseCount{"1/24/20", statistics{6, 4}},
			CaseCount{"1/25/20", statistics{10, 5}},
			CaseCount{"1/26/20", statistics{15, 8}},
		}},
		CountryCaseCounts{countryInformation{"United Kingdom", 55.3781, -3.4360000000000004}, []CaseCount{
			CaseCount{"1/23/20", statistics{6, 1}},
			CaseCount{"1/24/20", statistics{8, 3}},
			CaseCount{"1/25/20", statistics{9, 6}},
			CaseCount{"1/26/20", statistics{20, 6}},
		}},
	}
	verifyResultsCountryCaseCountsArr(result, expectedData, t)
}

func TestWorldTotal_AllDates(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetWorldCaseCounts("", "")
	if len(result) != 6 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 3)
	}
	expectedData := []CaseCount{
		CaseCount{"1/22/20", statistics{162, 35}},
		CaseCount{"1/23/20", statistics{1254, 198}},
		CaseCount{"1/24/20", statistics{2703, 257}},
		CaseCount{"1/25/20", statistics{3185, 328}},
		CaseCount{"1/26/20", statistics{3655, 376}},
		CaseCount{"1/27/20", statistics{3929, 456}},
	}
	verifyResultsCaseCountArr(result, expectedData, t)
}

func TestWorldTotal_QueryDates(t *testing.T) {
	caseCountsCache = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetWorldCaseCounts("1/23/20", "1/26/20")
	if len(result) != 4 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 3)
	}
	expectedData := []CaseCount{
		CaseCount{"1/23/20", statistics{1254, 198}},
		CaseCount{"1/24/20", statistics{2703, 257}},
		CaseCount{"1/25/20", statistics{3185, 328}},
		CaseCount{"1/26/20", statistics{3655, 376}},
	}
	verifyResultsCaseCountArr(result, expectedData, t)
}