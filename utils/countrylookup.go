package utils

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

const lookupURL = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/UID_ISO_FIPS_LookUp_Table.csv"

// AbbreviationToCountry : mapping of abbreviation to country name
var AbbreviationToCountry map[string]string

// CountryToAbbreviation : mapping of country name to abbreviation
var CountryToAbbreviation map[string]string

// StatePopulationLookup : mapping of country to state to population
var StatePopulationLookup map[string]map[string]int

var client HTTPClient

func init() {
	client = &http.Client{}
	getLookupData()
}

func getLookupData() {
	AbbreviationToCountry = make(map[string]string)
	CountryToAbbreviation = make(map[string]string)
	StatePopulationLookup = make(map[string]map[string]int)
	data, ok := ReadCSVFromURL(client, lookupURL)
	if !ok {
		log.Fatal("Unable to obtain lookup data, shutting down.")
	}
	populateAbbreviationCountryMaps(data[1:])
	populatePopulationMaps(data[1:])
}

func populateAbbreviationCountryMaps(data [][]string) {
	for _, row := range data {
		iso, state, country := row[1], row[6], row[7]
		if iso == "" || country == "" {
			continue
		}
		if _, ok := AbbreviationToCountry[iso]; !ok && state == "" {
			AbbreviationToCountry[iso] = country
			CountryToAbbreviation[country] = iso
		}
	}
}

func populatePopulationMaps(data [][]string) {
	for _, row := range data {
		state, country, population := row[6], row[7], row[11]
		if country == "" {
			continue
		}
		iso := CountryToAbbreviation[country]
		if _, ok := StatePopulationLookup[iso]; !ok {
			StatePopulationLookup[iso] = make(map[string]int)
		}
		popInt, err := strconv.Atoi(population)
		if err == nil {
			StatePopulationLookup[iso][state] = popInt
		}
	}
}

// GetCountryFromAbbreviation : get country name from iso code
func GetCountryFromAbbreviation(abbr string) (string, bool) {
	if _, ok := CountryToAbbreviation[abbr]; ok {
		// input is already a country
		return abbr, true
	}
	if country, ok := AbbreviationToCountry[strings.ToUpper(abbr)]; ok {
		return country, true
	}
	return "", false
}

// GetAbbreviationFromCountry : get iso code from country name
func GetAbbreviationFromCountry(country string) (string, bool) {
	if _, ok := AbbreviationToCountry[strings.ToUpper(country)]; ok {
		// input is already an iso
		return strings.ToUpper(country), true
	}
	if abbr, ok := CountryToAbbreviation[country]; ok {
		return abbr, true
	}
	return lowerCaseCountryLookup(country)
}

func lowerCaseCountryLookup(country string) (string, bool) {
	minEditDistance := -1
	closestMatch := ""
	for countryKey, iso := range CountryToAbbreviation {
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
