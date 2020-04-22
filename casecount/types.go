package casecount

type statistics struct {
	Confirmed int `json:"confirmed"`
	Deaths    int `json:"deaths"`
}

// CaseCount : contains stattics for given date
type CaseCount struct {
	Date string `json:"date"`
	statistics
}

// Location : point coordinates in the world map
type Location struct {
	Lat  float32 `json:"lat"`
	Long float32 `json:"long"`
}

// CaseCounts : contains information about the state,country and latitude longitude as well as the per day cumulative number of confirmed cases/deaths
type CaseCounts struct {
	Location
	Counts []CaseCount `json:"counts"`
}

// CaseCountsAggregated : contains the information about the state, country and the latitude/longitude as well as the number of confirmed cases/deaths
type CaseCountsAggregated struct {
	Location
	statistics
}
