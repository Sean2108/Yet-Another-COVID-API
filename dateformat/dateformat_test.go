package dateformat

import (
	"testing"
)

type testStruct struct {
	input      string
	expected   string
	expectedOk bool
}

func verifyResult(table testStruct, result string, ok bool, t *testing.T) {
	if table.expectedOk && !ok {
		t.Errorf("Format for %s was not found.", table.input)
	}
	if !table.expectedOk && ok {
		t.Errorf("Result should have been not ok for %s.", table.input)
	}
	if result != table.expected {
		t.Errorf("Result was incorrect, got: %s, want: %s.", result, table.expected)
	}
}

func TestGetCasesDateFormat(t *testing.T) {
	tables := []testStruct{
		{"2020-12-31", "12/31/20", true},
		{"12/31/20", "12/31/20", true},
		{"2020/12/31", "12/31/20", true},
		{"12-31-20", "12/31/20", true},
		{"2020/1/2", "1/2/20", true},
		{"2-1-20", "2/1/20", true},
		{"02-01-20", "2/1/20", true},
		{"20/1/2", "1/2/20", true},
		{"", "", true},
		{"2020-01-32", "", false},
		{"13/28/20", "", false},
	}

	for _, table := range tables {
		result, err := FormatDate(CasesDateFormat, table.input)
		verifyResult(table, result, err, t)
	}
}

func TestGetNewsDateFormat(t *testing.T) {
	tables := []testStruct{
		{"2020-01-31", "2020-01-31", true},
		{"1/31/20", "2020-01-31", true},
		{"", "", true},
		{"2020-01-32", "", false},
		{"13/28/20", "", false},
	}

	for _, table := range tables {
		result, ok := FormatDate(NewsDateFormat, table.input)
		verifyResult(table, result, ok, t)
	}
}
