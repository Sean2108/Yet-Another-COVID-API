package utils

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

const lookupURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/UID_ISO_FIPS_LookUp_Table.csv"

var abbreviationToCountry map[string]string
var countryToAbbreviation map[string]string
var statePopulationLookup map[string]map[string]int

var client HTTPClient

func init() {
	client = &http.Client{}
	getLookupData()
}

func getLookupData() {
	abbreviationToCountry = make(map[string]string)
	countryToAbbreviation = make(map[string]string)
	statePopulationLookup = make(map[string]map[string]int)
	data, err := ReadCSVFromURL(client, lookupURL)
	if err != nil {
		log.Fatal(err.Error())
	}
	parseLookupData(data[1:])
}

func parseLookupData(data [][]string) {
	for _, row := range data {
		iso, state, country, population := row[1], row[6], row[7], row[11]
		if iso == "" || country == "" {
			continue
		}
		if _, ok := abbreviationToCountry[iso]; !ok && state == "" {
			abbreviationToCountry[iso] = country
			countryToAbbreviation[country] = iso
		}
		if _, ok := statePopulationLookup[iso]; !ok {
			statePopulationLookup[iso] = make(map[string]int)
		}
		popInt, err := strconv.Atoi(population)
		if err == nil {
			statePopulationLookup[iso][state] = popInt
		}
	}
}

// GetCountryFromAbbreviation : get country name from iso code
func GetCountryFromAbbreviation(abbr string) (string, bool) {
	if _, ok := countryToAbbreviation[abbr]; ok {
		// input is already a country
		return abbr, true
	}
	if country, ok := abbreviationToCountry[strings.ToUpper(abbr)]; ok {
		return country, true
	}
	return "", false
}

// GetAbbreviationFromCountry : get iso code from country name
func GetAbbreviationFromCountry(country string) (string, bool) {
	if _, ok := abbreviationToCountry[strings.ToUpper(country)]; ok {
		// input is already an iso
		return strings.ToUpper(country), true
	}
	if abbr, ok := countryToAbbreviation[country]; ok {
		return abbr, true
	}
	return lowerCaseCountryLookup(country)
}

func lowerCaseCountryLookup(country string) (string, bool) {
	minEditDistance := -1
	closestMatch := ""
	for countryKey, iso := range countryToAbbreviation {
		if strings.ToLower(countryKey) == strings.ToLower(country) {
			return iso, true
		}
		if editDistance := editDistance([]rune(country), []rune(countryKey)); minEditDistance == -1 || editDistance < minEditDistance {
			minEditDistance = editDistance
			closestMatch = countryKey
		}
	}
	return closestMatch, false
}
