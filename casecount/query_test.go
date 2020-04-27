package casecount

import (
	"testing"
)

func getTestCaseCounts() map[string]CountryWithStates {
	result := map[string]CountryWithStates{
		"CN": CountryWithStates{
			Name: "China",
			States: map[string]CaseCounts{
				"Beijing": CaseCounts{
					LocationAndPopulation{40.1824, 116.4142, 50000},
					[]CaseCount{
						CaseCount{"1/22/20", statistics{50, 10, 0}},
						CaseCount{"1/23/20", statistics{200, 87, 10}},
						CaseCount{"1/24/20", statistics{800, 125, 30}},
						CaseCount{"1/25/20", statistics{1020, 142, 50}},
						CaseCount{"1/26/20", statistics{1110, 145, 60}},
						CaseCount{"1/27/20", statistics{1235, 152, 90}},
					},
				},
				"Hubei": CaseCounts{
					LocationAndPopulation{30.9756, 112.2707, 30000},
					[]CaseCount{
						CaseCount{"1/22/20", statistics{100, 20, 0}},
						CaseCount{"1/23/20", statistics{1000, 100, 50}},
						CaseCount{"1/24/20", statistics{1800, 105, 140}},
						CaseCount{"1/25/20", statistics{2020, 150, 240}},
						CaseCount{"1/26/20", statistics{2110, 175, 350}},
						CaseCount{"1/27/20", statistics{2111, 230, 460}},
					},
				},
				"Shanghai": CaseCounts{
					LocationAndPopulation{31.202, 121.4491, 40000},
					[]CaseCount{
						CaseCount{"1/22/20", statistics{10, 5, 0}},
						CaseCount{"1/23/20", statistics{45, 8, 2}},
						CaseCount{"1/24/20", statistics{89, 20, 4}},
						CaseCount{"1/25/20", statistics{126, 25, 5}},
						CaseCount{"1/26/20", statistics{400, 42, 7}},
						CaseCount{"1/27/20", statistics{532, 55, 10}},
					},
				},
			},
		},
		"SG": CountryWithStates{
			Name: "Singapore",
			States: map[string]CaseCounts{
				"": CaseCounts{
					LocationAndPopulation{1.2833, 103.8333, 6000},
					[]CaseCount{
						CaseCount{"1/22/20", statistics{1, 0, 0}},
						CaseCount{"1/23/20", statistics{3, 2, 0}},
						CaseCount{"1/24/20", statistics{6, 4, 1}},
						CaseCount{"1/25/20", statistics{10, 5, 2}},
						CaseCount{"1/26/20", statistics{15, 8, 4}},
						CaseCount{"1/27/20", statistics{23, 10, 6}},
					},
				},
			},
		},
		"GB": CountryWithStates{
			Name: "United Kingdom",
			States: map[string]CaseCounts{
				"London": CaseCounts{
					LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
					[]CaseCount{
						CaseCount{"1/22/20", statistics{1, 0, 0}},
						CaseCount{"1/23/20", statistics{6, 1, 0}},
						CaseCount{"1/24/20", statistics{8, 3, 0}},
						CaseCount{"1/25/20", statistics{9, 6, 2}},
						CaseCount{"1/26/20", statistics{20, 6, 5}},
						CaseCount{"1/27/20", statistics{28, 9, 10}},
					},
				},
			},
		},
	}
	return result
}

func getTestCaseCountsWithoutFirstAndLastDay() map[string]CountryWithStates {
	result := map[string]CountryWithStates{
		"CN": CountryWithStates{
			Name: "China",
			States: map[string]CaseCounts{
				"Beijing": CaseCounts{
					LocationAndPopulation{40.1824, 116.4142, 50000},
					[]CaseCount{
						CaseCount{"1/23/20", statistics{200, 87, 10}},
						CaseCount{"1/24/20", statistics{800, 125, 30}},
						CaseCount{"1/25/20", statistics{1020, 142, 50}},
						CaseCount{"1/26/20", statistics{1110, 145, 60}},
					},
				},
				"Hubei": CaseCounts{
					LocationAndPopulation{30.9756, 112.2707, 30000},
					[]CaseCount{
						CaseCount{"1/23/20", statistics{1000, 100, 50}},
						CaseCount{"1/24/20", statistics{1800, 105, 140}},
						CaseCount{"1/25/20", statistics{2020, 150, 240}},
						CaseCount{"1/26/20", statistics{2110, 175, 350}},
					},
				},
				"Shanghai": CaseCounts{
					LocationAndPopulation{31.202, 121.4491, 40000},
					[]CaseCount{
						CaseCount{"1/23/20", statistics{45, 8, 2}},
						CaseCount{"1/24/20", statistics{89, 20, 4}},
						CaseCount{"1/25/20", statistics{126, 25, 5}},
						CaseCount{"1/26/20", statistics{400, 42, 7}},
					},
				},
			},
		},
		"SG": CountryWithStates{
			Name: "Singapore",
			States: map[string]CaseCounts{
				"": CaseCounts{
					LocationAndPopulation{1.2833, 103.8333, 6000},
					[]CaseCount{
						CaseCount{"1/23/20", statistics{3, 2, 0}},
						CaseCount{"1/24/20", statistics{6, 4, 1}},
						CaseCount{"1/25/20", statistics{10, 5, 2}},
						CaseCount{"1/26/20", statistics{15, 8, 4}},
					},
				},
			},
		},
		"GB": CountryWithStates{
			Name: "United Kingdom",
			States: map[string]CaseCounts{
				"London": CaseCounts{
					LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
					[]CaseCount{
						CaseCount{"1/23/20", statistics{6, 1, 0}},
						CaseCount{"1/24/20", statistics{8, 3, 0}},
						CaseCount{"1/25/20", statistics{9, 6, 2}},
						CaseCount{"1/26/20", statistics{20, 6, 5}},
					},
				},
			},
		},
	}
	return result
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
	recoveredData := [][]string{
		{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"},
		{"Beijing", "China", "40.1824", "116.4142", "0", "10", "30", "50", "60", "90"},
		{"Hubei", "China", "30.9756", "112.2707", "0", "50", "140", "240", "350", "460"},
		{"Shanghai", "China", "31.202", "121.4491", "0", "2", "4", "5", "7", "10"},
		{"", "Singapore", "1.2833", "103.8333", "0", "0", "1", "2", "4", "6"},
		{"London", "United Kingdom", "55.3781", "-3.4360000000000004", "0", "0", "0", "2", "5", "10"},
		{"", "US", "37.0902", "-95.7129", "0", "0", "0", "50", "100", "150"},
	}
	headerRow := confirmedData[0]
	result := extractCaseCounts(headerRow, confirmedData, deathsData, recoveredData)
	expectedData := getTestCaseCounts()
	expectedData["US"] = CountryWithStates{
		Name: "US",
		States: map[string]CaseCounts{
			"": CaseCounts{
				LocationAndPopulation{37.0902, -95.7129, 300000},
				[]CaseCount{
					CaseCount{"1/22/20", statistics{0, 0, 0}},
					CaseCount{"1/23/20", statistics{0, 0, 0}},
					CaseCount{"1/24/20", statistics{0, 0, 0}},
					CaseCount{"1/25/20", statistics{0, 0, 50}},
					CaseCount{"1/26/20", statistics{0, 0, 100}},
					CaseCount{"1/27/20", statistics{0, 0, 150}},
				},
			},
		},
	}
	verifyResultsCaseCountsMap(result, expectedData, t)
}

func TestAggregateDataBetweenDates_AllDates(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result := stateAggregatedMap
	expectedData := map[string]CountryWithStatesAggregated{
		"CN": CountryWithStatesAggregated{
			Name: "China",
			States: map[string]CaseCountsAggregated{
				"Beijing": CaseCountsAggregated{
					LocationAndPopulation{40.1824, 116.4142, 50000},
					statistics{1235, 152, 90},
				},
				"Hubei": CaseCountsAggregated{
					LocationAndPopulation{30.9756, 112.2707, 30000},
					statistics{2111, 230, 460},
				},
				"Shanghai": CaseCountsAggregated{
					LocationAndPopulation{31.202, 121.4491, 40000},
					statistics{532, 55, 10},
				},
			},
		},
		"SG": CountryWithStatesAggregated{
			Name: "Singapore",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{1.2833, 103.8333, 6000},
					statistics{23, 10, 6},
				},
			},
		},
		"GB": CountryWithStatesAggregated{
			Name: "United Kingdom",
			States: map[string]CaseCountsAggregated{
				"London": CaseCountsAggregated{
					LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
					statistics{28, 9, 10},
				},
			},
		},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryDates(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/24/20", "1/26/20", "")
	expectedData := map[string]CountryWithStatesAggregated{
		"CN": CountryWithStatesAggregated{
			Name: "China",
			States: map[string]CaseCountsAggregated{
				"Beijing": CaseCountsAggregated{
					LocationAndPopulation{40.1824, 116.4142, 50000},
					statistics{910, 58, 50},
				},
				"Hubei": CaseCountsAggregated{
					LocationAndPopulation{30.9756, 112.2707, 30000},
					statistics{1110, 75, 300},
				},
				"Shanghai": CaseCountsAggregated{
					LocationAndPopulation{31.202, 121.4491, 40000},
					statistics{355, 34, 5},
				},
			},
		},
		"SG": CountryWithStatesAggregated{
			Name: "Singapore",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{1.2833, 103.8333, 6000},
					statistics{12, 6, 4},
				},
			},
		},
		"GB": CountryWithStatesAggregated{
			Name: "United Kingdom",
			States: map[string]CaseCountsAggregated{
				"London": CaseCountsAggregated{
					LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
					statistics{14, 5, 5},
				},
			},
		},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryDatesBeforeValidRange(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/20/20", "1/21/20", "")
	if len(result) != 0 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 0)
	}
}

func TestAggregateDataBetweenDates_QueryDatesAfterValidRange(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/28/20", "1/29/20", "")
	if len(result) != 0 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 0)
	}
}

func TestAggregateDataBetweenDates_QueryDatesBeforeAndAfter_ShouldReturnAll(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/21/20", "1/28/20", "")
	expectedData := map[string]CountryWithStatesAggregated{
		"CN": CountryWithStatesAggregated{
			Name: "China",
			States: map[string]CaseCountsAggregated{
				"Beijing": CaseCountsAggregated{
					LocationAndPopulation{40.1824, 116.4142, 50000},
					statistics{1235, 152, 90},
				},
				"Hubei": CaseCountsAggregated{
					LocationAndPopulation{30.9756, 112.2707, 30000},
					statistics{2111, 230, 460},
				},
				"Shanghai": CaseCountsAggregated{
					LocationAndPopulation{31.202, 121.4491, 40000},
					statistics{532, 55, 10},
				},
			},
		},
		"SG": CountryWithStatesAggregated{
			Name: "Singapore",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{1.2833, 103.8333, 6000},
					statistics{23, 10, 6},
				},
			},
		},
		"GB": CountryWithStatesAggregated{
			Name: "United Kingdom",
			States: map[string]CaseCountsAggregated{
				"London": CaseCountsAggregated{
					LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
					statistics{28, 9, 10},
				},
			},
		},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryFromDateAfterToDate(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	_, err := aggregateDataBetweenDates("1/24/20", "1/23/20", "CN")
	if err == nil {
		t.Error("Error message should be returned.")
	}
}

func TestAggregateDataBetweenDates_QueryCountry(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("", "", "SG")
	expectedData := map[string]CountryWithStatesAggregated{
		"SG": CountryWithStatesAggregated{
			Name: "Singapore",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{1.2833, 103.8333, 6000},
					statistics{23, 10, 6},
				},
			},
		},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryDates_FromIsOutOfRange(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/21/20", "1/26/20", "")
	expectedData := map[string]CountryWithStatesAggregated{
		"CN": CountryWithStatesAggregated{
			Name: "China",
			States: map[string]CaseCountsAggregated{
				"Beijing": CaseCountsAggregated{
					LocationAndPopulation{40.1824, 116.4142, 50000},
					statistics{1110, 145, 60},
				},
				"Hubei": CaseCountsAggregated{
					LocationAndPopulation{30.9756, 112.2707, 30000},
					statistics{2110, 175, 350},
				},
				"Shanghai": CaseCountsAggregated{
					LocationAndPopulation{31.202, 121.4491, 40000},
					statistics{400, 42, 7},
				},
			},
		},
		"SG": CountryWithStatesAggregated{
			Name: "Singapore",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{1.2833, 103.8333, 6000},
					statistics{15, 8, 4},
				},
			},
		},
		"GB": CountryWithStatesAggregated{
			Name: "United Kingdom",
			States: map[string]CaseCountsAggregated{
				"London": CaseCountsAggregated{
					LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
					statistics{20, 6, 5},
				},
			},
		},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryDates_ToIsOutOfRange(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/24/20", "1/28/20", "")
	expectedData := map[string]CountryWithStatesAggregated{
		"CN": CountryWithStatesAggregated{
			Name: "China",
			States: map[string]CaseCountsAggregated{
				"Beijing": CaseCountsAggregated{
					LocationAndPopulation{40.1824, 116.4142, 50000},
					statistics{1035, 65, 80},
				},
				"Hubei": CaseCountsAggregated{
					LocationAndPopulation{30.9756, 112.2707, 30000},
					statistics{1111, 130, 410},
				},
				"Shanghai": CaseCountsAggregated{
					LocationAndPopulation{31.202, 121.4491, 40000},
					statistics{487, 47, 8},
				},
			},
		},
		"SG": CountryWithStatesAggregated{
			Name: "Singapore",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{1.2833, 103.8333, 6000},
					statistics{20, 8, 6},
				},
			},
		},
		"GB": CountryWithStatesAggregated{
			Name: "United Kingdom",
			States: map[string]CaseCountsAggregated{
				"London": CaseCountsAggregated{
					LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
					statistics{22, 8, 10},
				},
			},
		},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataBetweenDates_QueryDates_FromAndToBothOutOfRange(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := aggregateDataBetweenDates("1/21/20", "1/28/20", "")
	expectedData := map[string]CountryWithStatesAggregated{
		"CN": CountryWithStatesAggregated{
			Name: "China",
			States: map[string]CaseCountsAggregated{
				"Beijing": CaseCountsAggregated{
					LocationAndPopulation{40.1824, 116.4142, 50000},
					statistics{1235, 152, 90},
				},
				"Hubei": CaseCountsAggregated{
					LocationAndPopulation{30.9756, 112.2707, 30000},
					statistics{2111, 230, 460},
				},
				"Shanghai": CaseCountsAggregated{
					LocationAndPopulation{31.202, 121.4491, 40000},
					statistics{532, 55, 10},
				},
			},
		},
		"SG": CountryWithStatesAggregated{
			Name: "Singapore",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{1.2833, 103.8333, 6000},
					statistics{23, 10, 6},
				},
			},
		},
		"GB": CountryWithStatesAggregated{
			Name: "United Kingdom",
			States: map[string]CaseCountsAggregated{
				"London": CaseCountsAggregated{
					LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
					statistics{28, 9, 10},
				},
			},
		},
	}
	verifyResultsCaseCountsAgg(result, expectedData, t)
}

func TestAggregateCountryDataFromStatesAggregate_AllDates(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result := countryAggregatedMap
	expectedData := map[string]CountryAggregated{
		"CN": CountryAggregated{
			"China",
			CaseCountsAggregated{
				LocationAndPopulation{(40.1824 + 30.9756 + 31.202) / 3.0, (116.4142 + 112.2707 + 121.4491) / 3.0, 120000},
				statistics{3878, 437, 560},
			},
		},
		"SG": CountryAggregated{
			"Singapore",
			CaseCountsAggregated{
				LocationAndPopulation{1.2833, 103.8333, 6000},
				statistics{23, 10, 6},
			},
		},
		"GB": CountryAggregated{
			"United Kingdom",
			CaseCountsAggregated{
				LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
				statistics{28, 9, 10},
			},
		},
	}
	verifyResultsCountryCaseCountsAgg(result, expectedData, t)
}

func TestAggregateCountryDataFromStatesAggregate_QueryDates(t *testing.T) {
	input := map[string]CountryWithStatesAggregated{
		"CN": CountryWithStatesAggregated{
			Name: "China",
			States: map[string]CaseCountsAggregated{
				"Beijing": CaseCountsAggregated{
					LocationAndPopulation{40.1824, 116.4142, 50000},
					statistics{910, 58, 50},
				},
				"Hubei": CaseCountsAggregated{
					LocationAndPopulation{30.9756, 112.2707, 30000},
					statistics{1110, 75, 300},
				},
				"Shanghai": CaseCountsAggregated{
					LocationAndPopulation{31.202, 121.4491, 40000},
					statistics{355, 34, 5},
				},
			},
		},
		"SG": CountryWithStatesAggregated{
			Name: "Singapore",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{1.2833, 103.8333, 6000},
					statistics{12, 6, 4},
				},
			},
		},
		"GB": CountryWithStatesAggregated{
			Name: "United Kingdom",
			States: map[string]CaseCountsAggregated{
				"": CaseCountsAggregated{
					LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
					statistics{14, 5, 5},
				},
			},
		},
	}
	result := aggregateCountryDataFromStatesAggregate(input)
	expectedData := map[string]CountryAggregated{
		"CN": CountryAggregated{
			"China",
			CaseCountsAggregated{
				LocationAndPopulation{(40.1824 + 30.9756 + 31.202) / 3.0, (116.4142 + 112.2707 + 121.4491) / 3.0, 120000},
				statistics{2375, 167, 355},
			},
		},
		"SG": CountryAggregated{
			"Singapore",
			CaseCountsAggregated{
				LocationAndPopulation{1.2833, 103.8333, 6000},
				statistics{12, 6, 4},
			},
		},
		"GB": CountryAggregated{
			"United Kingdom",
			CaseCountsAggregated{
				LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
				statistics{14, 5, 5},
			},
		},
	}
	verifyResultsCountryCaseCountsAgg(result, expectedData, t)
}

func TestAggregateDataPerDay_AllDates(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCaseCountsWithDayData("", "", "")
	expectedData := caseCountsMap
	verifyResultsCaseCountsMap(result, expectedData, t)
}

func TestAggregateDataPerDay_QueryDates(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCaseCountsWithDayData("1/23/20", "1/26/20", "")
	expectedData := getTestCaseCountsWithoutFirstAndLastDay()
	verifyResultsCaseCountsMap(result, expectedData, t)
}

func TestAggregateDataPerDay_BeforeAndAfterShouldReturnAll(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCaseCountsWithDayData("1/21/20", "1/28/20", "")
	expectedData := caseCountsMap
	verifyResultsCaseCountsMap(result, expectedData, t)
}

func TestAggregateDataPerDay_CountryQuery(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCaseCountsWithDayData("", "", "CN")
	expectedData := getTestCaseCounts()["CN"]
	verifyResultsCaseCountsMap(result, map[string]CountryWithStates{"CN": expectedData}, t)
}

func TestAggregateDataPerDay_QueryFromDateAfterToDate(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	_, err := GetCaseCountsWithDayData("1/24/20", "1/23/20", "CN")
	if err == nil {
		t.Error("Error message should be returned.")
	}
}

func TestCountryAggregateDataPerDay_AllDates(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCountryCaseCountsWithDayData("", "", "")
	expectedData := countryCaseCountsMap
	verifyResultsCountryCaseCountsMap(result, expectedData, t)
}

func TestCountryAggregateDataPerDay_QueryDates(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetCountryCaseCountsWithDayData("1/23/20", "1/26/20", "")
	expectedData := map[string]Country{
		"CN": Country{
			"China",
			CaseCounts{
				LocationAndPopulation{(40.1824 + 30.9756 + 31.202) / 3.0, (116.4142 + 112.2707 + 121.4491) / 3.0, 120000},
				[]CaseCount{
					CaseCount{"1/23/20", statistics{1245, 195, 62}},
					CaseCount{"1/24/20", statistics{2689, 250, 174}},
					CaseCount{"1/25/20", statistics{3166, 317, 295}},
					CaseCount{"1/26/20", statistics{3620, 362, 417}},
				},
			},
		},
		"SG": Country{
			"Singapore",
			CaseCounts{
				LocationAndPopulation{1.2833, 103.8333, 6000},
				[]CaseCount{
					CaseCount{"1/23/20", statistics{3, 2, 0}},
					CaseCount{"1/24/20", statistics{6, 4, 1}},
					CaseCount{"1/25/20", statistics{10, 5, 2}},
					CaseCount{"1/26/20", statistics{15, 8, 4}},
				},
			},
		},
		"GB": Country{
			"United Kingdom",
			CaseCounts{
				LocationAndPopulation{55.3781, -3.4360000000000004, 7000},
				[]CaseCount{
					CaseCount{"1/23/20", statistics{6, 1, 0}},
					CaseCount{"1/24/20", statistics{8, 3, 0}},
					CaseCount{"1/25/20", statistics{9, 6, 2}},
					CaseCount{"1/26/20", statistics{20, 6, 5}},
				},
			},
		},
	}
	verifyResultsCountryCaseCountsMap(result, expectedData, t)
}

func TestWorldTotal_AllDates(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetWorldCaseCounts("", "")
	if len(result) != 6 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 3)
	}
	expectedData := []CaseCount{
		CaseCount{"1/22/20", statistics{162, 35, 0}},
		CaseCount{"1/23/20", statistics{1254, 198, 62}},
		CaseCount{"1/24/20", statistics{2703, 257, 175}},
		CaseCount{"1/25/20", statistics{3185, 328, 299}},
		CaseCount{"1/26/20", statistics{3655, 376, 426}},
		CaseCount{"1/27/20", statistics{3929, 456, 576}},
	}
	verifyResultsCaseCountArr(result, expectedData, t)
}

func TestWorldTotal_QueryDates(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	result, _ := GetWorldCaseCounts("1/23/20", "1/26/20")
	if len(result) != 4 {
		t.Errorf("Length of results is incorrect, got: %d, want %d.", len(result), 3)
	}
	expectedData := []CaseCount{
		CaseCount{"1/23/20", statistics{1254, 198, 62}},
		CaseCount{"1/24/20", statistics{2703, 257, 175}},
		CaseCount{"1/25/20", statistics{3185, 328, 299}},
		CaseCount{"1/26/20", statistics{3655, 376, 426}},
	}
	verifyResultsCaseCountArr(result, expectedData, t)
}

func TestWorldTotal_QueryFromDateAfterToDate(t *testing.T) {
	caseCountsMap = getTestCaseCounts()
	setDateBoundariesAndAllAggregatedData([]string{"Province/State", "Country/Region", "Lat", "Long", "1/22/20", "1/23/20", "1/24/20", "1/25/20", "1/26/20", "1/27/20"})
	_, err := GetWorldCaseCounts("1/24/20", "1/23/20")
	if err == nil {
		t.Error("Error message should be returned.")
	}
}
