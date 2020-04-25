package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type mockClient struct{}

func (m *mockClient) Get(url string) (*http.Response, error) {
	csvStr := "UID,iso2,iso3,code3,FIPS,Admin2,Province_State,Country_Region,Lat,Long_,Combined_Key,Population\n4,AF,AFG,4,,,,Afghanistan,33.93911,67.709953,Afghanistan,38928341\n8,AL,ALB,8,,,,Albania,41.1533,20.1683,Albania,2877800\n12,DZ,DZA,12,,,,Algeria,28.0339,1.6596,Algeria,43851043"
	r := ioutil.NopCloser(bytes.NewReader([]byte(csvStr)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

func TestGetCountryFromAbbreviation(t *testing.T) {
	client = &mockClient{}
	getLookupData()

	tables := []struct {
		iso     string
		country string
		ok      bool
	}{
		{"aF", "Afghanistan", true},
		{"AF", "Afghanistan", true},
		{"AL", "Albania", true},
		{"DZ", "Algeria", true},
		{"Unknown", "", false},
	}
	for _, table := range tables {
		country, ok := GetCountryFromAbbreviation(table.iso)
		if table.ok != ok {
			t.Errorf("ok is incorrect, got: %t, want: %t.", ok, table.ok)
		}
		if table.country != country {
			t.Errorf("country is incorrect, got: %s, want: %s.", country, table.country)
		}
	}
}

func TestGetAbbreviationFromCountry(t *testing.T) {
	client = &mockClient{}
	getLookupData()

	tables := []struct {
		iso     string
		country string
		ok      bool
	}{
		{"AF", "AF", true},
		{"AF", "aF", true},
		{"AL", "Albania", true},
		{"DZ", "Algeria", true},
		{"AL", "aLbaNiA", true},
		{"Albania", "aaLbaNiA", false},
	}
	for _, table := range tables {
		iso, ok := GetAbbreviationFromCountry(table.country)
		if table.ok != ok {
			t.Errorf("ok is incorrect, got: %t, want: %t.", ok, table.ok)
		}
		if table.iso != iso {
			t.Errorf("country is incorrect, got: %s, want: %s.", iso, table.iso)
		}
	}
}
