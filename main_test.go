package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"yet-another-covid-map-api/casecount"
)

func getData(URL string, t *testing.T) map[string]map[string]casecount.CaseCountsAggregated {
	resp, err := http.Get(URL)
	if err != nil {
		t.Errorf("Error when retrieving from cases endpoint: %s.", err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error when reading body of response: %s.", err.Error())
	}
	var aggregatedData map[string]map[string]casecount.CaseCountsAggregated
	json.Unmarshal(body, &aggregatedData)
	return aggregatedData
}

func verifyData(country string, state string, allItem casecount.CaseCountsAggregated, queriedItem casecount.CaseCountsAggregated, t *testing.T) {
	if allItem.Confirmed < queriedItem.Confirmed {
		t.Errorf("AllAggregateData has fewer confirmed than queriedAggregateData, country: %s, state: %s, all: %d, queried: %d", country, state, allItem.Confirmed, queriedItem.Confirmed)
	}
	if allItem.Deaths < queriedItem.Deaths {
		t.Errorf("AllAggregateData has fewer deaths than queriedAggregateData, country: %s, state: %s, all: %d, queried: %d", country, state, allItem.Deaths, queriedItem.Deaths)
	}
}

func TestCasesEndpoint(t *testing.T) {
	go main()
	time.Sleep(5 * time.Second)
	allAggregatedData := getData("http://localhost:"+port+"/cases", t)
	if len(allAggregatedData) == 0 {
		t.Errorf("Response has no items.")
	}
	queriedAggregatedData := getData("http://localhost:"+port+"/cases?from=3/15/20&to=3/17/20", t)
	if len(queriedAggregatedData) == 0 {
		t.Errorf("Response has no items.")
	}

	if len(allAggregatedData) != len(queriedAggregatedData) {
		t.Errorf("Reponses have different lengths, all: %d, queried: %d", len(allAggregatedData), len(queriedAggregatedData))
	}
	for country, countryInfo := range allAggregatedData {
		for state, stateInfo := range countryInfo {
			verifyData(country, state, stateInfo, queriedAggregatedData[country][state], t)
		}
	}
}
