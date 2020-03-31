package casecount

type ByCountryAndStateForCaseCounts []caseCounts

func (a ByCountryAndStateForCaseCounts) Len() int {
	return len(a)
}

func (a ByCountryAndStateForCaseCounts) Less(i, j int) bool {
	if a[i].Country == a[j].Country {
		return a[i].State < a[j].State
	}
	return a[i].Country < a[j].Country
}

func (a ByCountryAndStateForCaseCounts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ByCountryAndStateAgg []CaseCountsAggregated

func (a ByCountryAndStateAgg) Len() int {
	return len(a)
}

func (a ByCountryAndStateAgg) Less(i, j int) bool {
	if a[i].Country == a[j].Country {
		return a[i].State < a[j].State
	}
	return a[i].Country < a[j].Country
}

func (a ByCountryAndStateAgg) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ByCountryAgg []CountryCaseCountsAggregated

func (a ByCountryAgg) Len() int {
	return len(a)
}

func (a ByCountryAgg) Less(i, j int) bool {
	return a[i].Country < a[j].Country
}

func (a ByCountryAgg) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
