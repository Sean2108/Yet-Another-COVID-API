package utils

import "strings"

var abbreviationToCountry map[string]string
var countryToAbbreviation map[string]string

func init() {
	abbreviationToCountry = map[string]string{
		"ae": "United Arab Emirates",
		"ar": "Argentina",
		"at": "Austria",
		"au": "Australia",
		"be": "Belgium",
		"bg": "Bulgaria",
		"br": "Brazil",
		"ca": "Canada",
		"ch": "Switzerland",
		"cn": "China",
		"co": "Colombia",
		"cu": "Cuba",
		"cz": "Czechia",
		"de": "Germany",
		"eg": "Egypt",
		"fr": "France",
		"gb": "United Kingdom",
		"gr": "Greece",
		"hk": "Hong Kong",
		"hu": "Hungary",
		"id": "Indonesia",
		"ie": "Ireland",
		"il": "Israel",
		"in": "India",
		"it": "Italy",
		"jp": "Japan",
		"kr": "Korea, South",
		"lt": "Lithuania",
		"lv": "Latvia",
		"ma": "Morocco",
		"mx": "Mexico",
		"my": "Malaysia",
		"ng": "Nigeria",
		"nl": "Netherlands",
		"no": "Norway",
		"nz": "New Zealand",
		"ph": "Philippines",
		"pl": "Poland",
		"pt": "Portugal",
		"ro": "Romania",
		"rs": "Serbia",
		"ru": "Russia",
		"sa": "Saudi Arabia",
		"se": "Sweden",
		"sg": "Singapore",
		"si": "Slovenia",
		"sk": "Slovakia",
		"th": "Thailand",
		"tr": "Turkey",
		"tw": "Taiwan",
		"ua": "Ukraine",
		"us": "US",
		"ve": "Venezuela",
		"za": "South Africa",
	}
	countryToAbbreviation = reverseMap(abbreviationToCountry)
}

func reverseMap(m map[string]string) map[string]string {
	reversedMap := make(map[string]string)
	for k, v := range m {
		reversedMap[v] = k
	}
	return reversedMap
}

func GetCountryFromAbbreviation(abbr string) (string, bool) {
	if country, ok := abbreviationToCountry[strings.ToLower(abbr)]; ok {
		return country, true
	}
	return "", false
}

func GetAbbreviationFromCountry(country string) (string, bool) {
	if abbr, ok := countryToAbbreviation[country]; ok {
		return abbr, true
	}
	return "", false
}
