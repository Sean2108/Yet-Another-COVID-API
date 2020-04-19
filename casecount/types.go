package casecount

type statistics struct {
	Confirmed int
	Deaths    int
}

type CaseCount struct {
	Date string
	statistics
}

type stateInformation struct {
	State   string
	Country string
	Lat     float32
	Long    float32
}

type countryInformation struct {
	Country string
	Lat     float32
	Long    float32
}

// CaseCounts : contains information about the state,country and latitude longitude as well as the per day cumulative number of confirmed cases/deaths
type CaseCounts struct {
	stateInformation
	Counts []CaseCount
}

// CountryCaseCounts : contains information about the state,country and latitude longitude as well as the per day cumulative number of confirmed cases/deaths
type CountryCaseCounts struct {
	countryInformation
	Counts []CaseCount
}

// CaseCountsAggregated : contains the information about the state, country and the latitude/longitude as well as the number of confirmed cases/deaths
type CaseCountsAggregated struct {
	stateInformation
	statistics
}

// CountryCaseCountsAggregated : contains the information about the country and the latitude/longitude as well as the number of confirmed cases/deaths
type CountryCaseCountsAggregated struct {
	countryInformation
	statistics
}
