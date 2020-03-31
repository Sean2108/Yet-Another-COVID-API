package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
	"testing"

	"yet-another-covid-map-api/casecount"
)

func getData(URL string, t *testing.T) []casecount.CaseCountsAggregated {
	resp, err := http.Get(URL)
	if err != nil {
		t.Errorf("Error when retrieving from cases endpoint: %s.", err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error when reading body of response: %s.", err.Error())
	}
	var aggregatedData []casecount.CaseCountsAggregated
	json.Unmarshal(body, &aggregatedData)
	return aggregatedData
}

func verifyData(allItem casecount.CaseCountsAggregated, queriedItem casecount.CaseCountsAggregated, t *testing.T) {
	if allItem.Country != queriedItem.Country {
		t.Errorf("AllAggregateData has a different country than queriedAggregateData, all: %s, queried: %s", allItem.Country, queriedItem.Country)
	}
	if allItem.State != queriedItem.State {
		t.Errorf("AllAggregateData has a different state than queriedAggregateData, country: %s, all: %s, queried: %s", allItem.Country, allItem.State, queriedItem.State)
	}
	if allItem.Confirmed < queriedItem.Confirmed {
		t.Errorf("AllAggregateData has fewer confirmed than queriedAggregateData, country: %s, state: %s, all: %d, queried: %d", allItem.Country, allItem.State, allItem.Confirmed, queriedItem.Confirmed)
	}
	if allItem.Deaths < queriedItem.Deaths {
		t.Errorf("AllAggregateData has fewer deaths than queriedAggregateData, country: %s, state: %s, all: %d, queried: %d", allItem.Country, allItem.State, allItem.Deaths, queriedItem.Deaths)
	}
}

func TestCasesEndpoint(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	go main()
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
	sort.Sort(casecount.ByCountryAndStateAgg(allAggregatedData))
	sort.Sort(casecount.ByCountryAndStateAgg(queriedAggregatedData))
	for i := 0; i < len(allAggregatedData); i++ {
		verifyData(allAggregatedData[i], queriedAggregatedData[i], t)
	}
}
