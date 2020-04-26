package casecount

type statistics struct {
	Confirmed int `json:"confirmed"`
	Deaths    int `json:"deaths"`
	Recovered int `json:"recovered"`
}

// CaseCount : contains stattics for given date
type CaseCount struct {
	Date string `json:"date"`
	statistics
}

// LocationAndPopulation : point coordinates in the world map and population of state/country
type LocationAndPopulation struct {
	Lat        float32 `json:"lat"`
	Long       float32 `json:"long"`
	Population int     `json:"population"`
}

// CaseCounts : contains information about the state,country and latitude longitude as well as the per day cumulative number of confirmed cases/deaths
type CaseCounts struct {
	LocationAndPopulation
	Counts []CaseCount `json:"counts"`
}

// CaseCountsAggregated : contains the information about the state, country and the latitude/longitude as well as the number of confirmed cases/deaths
type CaseCountsAggregated struct {
	LocationAndPopulation
	statistics
}

// CountryWithStates : contains name and state information of the country with detailed states information
type CountryWithStates struct {
	Name   string                `json:"country"`
	States map[string]CaseCounts `json:"states"`
}

// Country : contains name and information of the country
type Country struct {
	Name string `json:"country"`
	CaseCounts
}

// CountryWithStatesAggregated : contains name and aggregated state information of the country with detailed states information
type CountryWithStatesAggregated struct {
	Name   string                          `json:"country"`
	States map[string]CaseCountsAggregated `json:"states"`
}

// CountryAggregated : contains name and aggregated information of the country
type CountryAggregated struct {
	Name string `json:"country"`
	CaseCountsAggregated
}
