package requests

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"yet-another-covid-map-api/casecount"
	"yet-another-covid-map-api/dateformat"
)

var (
	fakeResponse []byte
	testFnCalled bool
)

type fakeWriter struct{}

func (w *fakeWriter) Header() http.Header {
	return http.Header{}
}

func (w *fakeWriter) Write(response []byte) (int, error) {
	fakeResponse = response
	return 0, nil
}

func (w *fakeWriter) WriteHeader(statusCode int) {
	return
}

func callTestFn(from string, to string, country string, aggregateCountries bool, perDay bool, worldTotal bool) ([]byte, error, error) {
	testFnCalled = true
	return []byte("response"), nil, nil
}

func TestParseUrlQuery(t *testing.T) {
	tables := []struct {
		rawurl             string
		from               string
		to                 string
		country            string
		aggregateCountries bool
		perDay             bool
		worldTotal         bool
		errorString        string
	}{
		{"http://localhost:8080/cases", "", "", "", false, false, false, ""},
		{"http://localhost:8080/cases?from=&to=1/1/20", "", "1/1/20", "", false, false, false, ""},
		{"http://localhost:8080/cases?from=1/1/20&to=1/2/20", "1/1/20", "1/2/20", "", false, false, false, ""},
		{"http://localhost:8080/cases?from=&to=1/1/20", "", "1/1/20", "", false, false, false, ""},
		{"http://localhost:8080/cases?from=&to=&country=CN", "", "", "CN", false, false, false, ""},
		{"http://localhost:8080/cases?from=&to=&country=gb", "", "", "GB", false, false, false, ""},
		{"http://localhost:8080/cases?country=United Kingdom", "", "", "GB", false, false, false, ""},
		{"http://localhost:8080/cases?aggregateCountries=true&country=sg", "", "", "SG", true, false, false, ""},
		{"http://localhost:8080/cases?aggregatecountries=true&country=sg", "", "", "SG", true, false, false, ""},
		{"http://localhost:8080/cases?aggregateCountries=tru&country=Singapore", "", "", "SG", false, false, false, ""},
		{"http://localhost:8080/cases?aggregateCountries=tru&country=Sngapore", "", "", "", false, false, false, "Singapore"},
		{"http://localhost:8080/cases?from=1/1/20&to=1/2/20&country=Singapore&aggregateCountries=true&perDay=false", "1/1/20", "1/2/20", "SG", true, false, false, ""},
		{"http://localhost:8080/cases?from=1/1/20&to=1/2/20&country=Singapore&aggregateCountries=false&perDay=true", "1/1/20", "1/2/20", "SG", false, true, false, ""},
		{"http://localhost:8080/cases?from=1/32/20&to=1/2/20&country=Singapore&aggregateCountries=true", "", "", "", false, false, false, "Date format"},
		{"http://localhost:8080/cases?from=1/1/20&to=1/32/20&country=Singapore&aggregateCountries=true", "", "", "", false, false, false, "Date format"},
		{"http://localhost:8080/cases?from=1/1/20&to=1/30/20&worldTotal=true", "1/1/20", "1/30/20", "", false, false, true, ""},
		{"http://localhost:8080/cases?from=1/1/20&to=1/30/20&worldTotal=false", "1/1/20", "1/30/20", "", false, false, false, ""},
	}

	for _, table := range tables {
		url, _ := url.Parse(table.rawurl)
		from, to, country, aggregateCountries, perDay, worldTotal, err := parseURL(url, dateformat.CasesDateFormat)
		if from != table.from {
			t.Errorf("from result of parseURL was incorrect for %s, got: %s, want: %s.", table.rawurl, from, table.from)
		}
		if to != table.to {
			t.Errorf("to result of parseURL was incorrect for %s, got: %s, want: %s.", table.rawurl, to, table.to)
		}
		if country != table.country {
			t.Errorf("country result of parseURL was incorrect for %s, got: %s, want: %s.", table.rawurl, country, table.country)
		}
		if aggregateCountries != table.aggregateCountries {
			t.Errorf("aggregateCountries result of parseURL was incorrect for %s, got: %t, want: %t.", table.rawurl, aggregateCountries, table.aggregateCountries)
		}
		if perDay != table.perDay {
			t.Errorf("perDay result of parseURL was incorrect for %s, got: %t, want: %t.", table.rawurl, perDay, table.perDay)
		}
		if worldTotal != table.worldTotal {
			t.Errorf("worldTotal result of parseURL was incorrect for %s, got: %t, want: %t.", table.rawurl, worldTotal, table.worldTotal)
		}
		if table.errorString != "" && !strings.Contains(err.Error(), table.errorString) {
			t.Errorf("error thrown for parseURL was incorrect for %s, got: %s, want error containing: %s.", table.rawurl, err.Error(), table.errorString)
		}
	}
}

func TestGetCaseCountsResponse_PerDay(t *testing.T) {

	tables := []struct {
		country            string
		aggregateCountries bool
		perDay             bool
		worldTotal         bool
	}{
		{"", false, true, false},
		{"SG", false, true, false},
		{"", true, false, false},
		{"SG", true, false, false},
		{"", false, false, false},
		{"", true, true, false},
		{"SG", false, false, false},
		{"", false, false, true},
	}

	for _, table := range tables {
		casecount.UpdateCaseCounts()
		response, err, caseCountErr := getCaseCountsResponse("", "", table.country, table.aggregateCountries, table.perDay, table.worldTotal)
		if len(response) < 3 {
			t.Errorf("Response should not be empty, got length: %d, want length: %s.", len(response), "more than 2")
		}
		if err != nil {
			t.Errorf("Err should be null, got: %s, want: nil.", err.Error())
		}
		if caseCountErr != nil {
			t.Errorf("caseCountErr should be null, got: %s, want: nil.", caseCountErr.Error())
		}
	}
}
func TestGetNewsForCountryResponse_PerDay(t *testing.T) {
	response, err, newsErr := getNewsForCountryResponse("", "", "Singapore", false, false, false)
	if len(response) < 3 {
		t.Errorf("Response should not be empty, got length: %d, want length: %s.", len(response), "more than 2")
	}
	if err != nil {
		t.Errorf("Err should be null, got: %s, want: nil.", err.Error())
	}
	if newsErr != nil {
		t.Errorf("caseCountErr should be null, got: %s, want: nil.", newsErr.Error())
	}
}

func TestGetResponse(t *testing.T) {
	testURLs := []string{
		"http://localhost:8080/cases?country=sg",
		"http://localhost:8080/cases?worldTotal=true",
		"http://localhost:8080/news?country=Singapore",
	}
	for _, testURL := range testURLs {
		fakeResponse = []byte("")
		testFnCalled = false
		inputURL, _ := url.Parse(testURL)
		getResponse(callTestFn, &fakeWriter{}, inputURL, true)
		if !testFnCalled {
			t.Error("callTestFn should have been called, but it was not.")
		}
		if string(fakeResponse) != "response" {
			t.Errorf("fakeResponse should have been modified, got: %s, want: response", fakeResponse)
		}
	}
}

func TestGetResponse_shouldFailWhenDateIsMalformed(t *testing.T) {
	fakeResponse = []byte("")
	testFnCalled = false
	inputURL, _ := url.Parse("http://localhost:8080/cases?from=3/32/20")
	getResponse(callTestFn, &fakeWriter{}, inputURL, false)
	if testFnCalled {
		t.Error("callTestFn should not have been called, but it was.")
	}
	if !strings.Contains(string(fakeResponse), "Date format is not recognised") {
		t.Errorf("fakeResponse did not contain the correct error message, got: %s, want: string containing message about date format not recognised", fakeResponse)
	}
}
