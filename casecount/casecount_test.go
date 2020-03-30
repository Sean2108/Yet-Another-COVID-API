package casecount

import (
	"testing"
	"time"
)

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

func TestGetStatisticsSum_EmptyInput(t *testing.T) {
	var input []caseCount
	confirmed, deaths := getStatisticsSum(input)
	if confirmed != 0 {
		t.Errorf("Confirmed was not 0, got: %d.", confirmed)
	}
	if deaths != 0 {
		t.Errorf("Deaths was not 0, got: %d.", deaths)
	}
}

func TestGetStatisticsSum(t *testing.T) {
	var input = []caseCount{
		caseCount{"a", statistics{2, 1}},
		caseCount{"b", statistics{3, 2}},
		caseCount{"c", statistics{4, 3}},
	}
	confirmed, deaths := getStatisticsSum(input)
	if confirmed != 9 {
		t.Errorf("Confirmed was not correct, got: %d, want %d.", confirmed, 9)
	}
	if deaths != 6 {
		t.Errorf("Deaths was not 0, got: %d, want %d.", deaths, 6)
	}
}
