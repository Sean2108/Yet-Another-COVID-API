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
		getAbbreviation    bool
		from               string
		to                 string
		country            string
		aggregateCountries bool
		perDay             bool
		worldTotal         bool
	}{
		{"http://localhost:8080/cases", false, "", "", "", false, false, false},
		{"http://localhost:8080/cases?from=&to=1/1/20", false, "", "1/1/20", "", false, false, false},
		{"http://localhost:8080/cases?from=1/1/20&to=1/2/20", false, "1/1/20", "1/2/20", "", false, false, false},
		{"http://localhost:8080/cases?from=&to=1/1/20", false, "", "1/1/20", "", false, false, false},
		{"http://localhost:8080/cases?from=&to=&country=CN", false, "", "", "China", false, false, false},
		{"http://localhost:8080/cases?from=&to=&country=gb", false, "", "", "United Kingdom", false, false, false},
		{"http://localhost:8080/cases?country=United Kingdom", false, "", "", "United Kingdom", false, false, false},
		{"http://localhost:8080/cases?aggregateCountries=true&country=sg", false, "", "", "Singapore", true, false, false},
		{"http://localhost:8080/cases?aggregatecountries=true&country=sg", false, "", "", "Singapore", true, false, false},
		{"http://localhost:8080/cases?aggregateCountries=tru&country=Singapore", true, "", "", "sg", false, false, false},
		{"http://localhost:8080/cases?aggregateCountries=tru&country=Sngapore", true, "", "", "Sngapore", false, false, false},
		{"http://localhost:8080/cases?from=1/1/20&to=1/2/20&country=Singapore&aggregateCountries=true&perDay=false", false, "1/1/20", "1/2/20", "Singapore", true, false, false},
		{"http://localhost:8080/cases?from=1/1/20&to=1/2/20&country=Singapore&aggregateCountries=false&perDay=true", false, "1/1/20", "1/2/20", "Singapore", false, true, false},
		{"http://localhost:8080/cases?from=1/32/20&to=1/2/20&country=Singapore&aggregateCountries=true", false, "", "", "", false, false, false},
		{"http://localhost:8080/cases?from=1/1/20&to=1/32/20&country=Singapore&aggregateCountries=true", false, "", "", "", false, false, false},
		{"http://localhost:8080/cases?from=1/1/20&to=1/30/20&worldTotal=true", false, "1/1/20", "1/30/20", "", false, false, true},
		{"http://localhost:8080/cases?from=1/1/20&to=1/30/20&worldTotal=false", false, "1/1/20", "1/30/20", "", false, false, false},
	}

	for _, table := range tables {
		url, _ := url.Parse(table.rawurl)
		from, to, country, aggregateCountries, perDay, worldTotal, _ := parseURL(url, table.getAbbreviation, dateformat.CasesDateFormat)
		if from != table.from {
			t.Errorf("result of parseURL was incorrect for %s, got: %s, want: %s.", table.rawurl, from, table.from)
		}
		if to != table.to {
			t.Errorf("result of parseURL was incorrect for %s, got: %s, want: %s.", table.rawurl, to, table.to)
		}
		if country != table.country {
			t.Errorf("result of parseURL was incorrect for %s, got: %s, want: %s.", table.rawurl, country, table.country)
		}
		if aggregateCountries != table.aggregateCountries {
			t.Errorf("result of parseURL was incorrect for %s, got: %t, want: %t.", table.rawurl, aggregateCountries, table.aggregateCountries)
		}
		if perDay != table.perDay {
			t.Errorf("result of parseURL was incorrect for %s, got: %t, want: %t.", table.rawurl, perDay, table.perDay)
		}
		if worldTotal != table.worldTotal {
			t.Errorf("result of parseURL was incorrect for %s, got: %t, want: %t.", table.rawurl, worldTotal, table.worldTotal)
		}
	}
}

func TestGetCaseCountsResponse_PerDay(t *testing.T) {

	tables := []struct {
		country            string
		aggregateCountries bool
		perDay             bool
		expectedSuccess    bool
	}{
		{"", false, true, true},
		{"Ssingapore", false, true, false},
		{"", true, false, true},
		{"Ssingapore", true, false, false},
		{"", false, false, true},
		{"Ssingapore", false, false, false},
	}

	for _, table := range tables {
		casecount.UpdateCaseCounts()
		response, err, caseCountErr := getCaseCountsResponse("", "", table.country, table.aggregateCountries, table.perDay, false)
		if table.expectedSuccess {
			if len(response) < 3 {
				t.Errorf("Response should not be empty, got length: %d, want length: %s.", len(response), "more than 2")
			}
			if err != nil {
				t.Errorf("Err should be null, got: %s, want: nil.", err.Error())
			}
			if caseCountErr != nil {
				t.Errorf("caseCountErr should be null, got: %s, want: nil.", caseCountErr.Error())
			}
		} else {
			if caseCountErr == nil || !strings.Contains(caseCountErr.Error(), "Singapore") {
				t.Errorf("caseCountErr is incorrect, got: %s.", caseCountErr.Error())
			}
		}
	}
}

func TestGetResponse(t *testing.T) {
	fakeResponse = []byte("")
	testFnCalled = false
	inputURL, _ := url.Parse("http://localhost:8080/cases?country=sg")
	getResponse(callTestFn, &fakeWriter{}, inputURL)
	if !testFnCalled {
		t.Error("callTestFn should have been called, but it was not.")
	}
	if string(fakeResponse) != "response" {
		t.Errorf("fakeResponse should have been modified, got: %s, want: response", fakeResponse)
	}
}

func TestGetResponse_shouldFailWhenDateIsMalformed(t *testing.T) {
	fakeResponse = []byte("")
	testFnCalled = false
	inputURL, _ := url.Parse("http://localhost:8080/cases?from=3/32/20")
	getResponse(callTestFn, &fakeWriter{}, inputURL)
	if testFnCalled {
		t.Error("callTestFn should not have been called, but it was.")
	}
	if !strings.Contains(string(fakeResponse), "Date format is not recognised") {
		t.Errorf("fakeResponse did not contain the correct error message, got: %s, want: string containing message about date format not recognised", fakeResponse)
	}
}
