package casecount

type statistics struct {
	Confirmed int
	Deaths    int
}

type CaseCount struct {
	Date string
	statistics
}

type Location struct {
	Lat  float32
	Long float32
}

// CaseCounts : contains information about the state,country and latitude longitude as well as the per day cumulative number of confirmed cases/deaths
type CaseCounts struct {
	Location
	Counts []CaseCount
}

// CaseCountsAggregated : contains the information about the state, country and the latitude/longitude as well as the number of confirmed cases/deaths
type CaseCountsAggregated struct {
	Location
	statistics
}
