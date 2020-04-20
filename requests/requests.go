package requests

import (
	"net/http"
)

// GetCaseCounts : logic when /cases endpoint is called. Returns all aggregated confirmed cases/death counts between from and to dates in the query
func GetCaseCounts(w http.ResponseWriter, r *http.Request) {
	getResponse(getCaseCountsResponse, w, r.URL, false)
}

// GetNewsForCountry : runs query to get all virus related news for a given country
func GetNewsForCountry(w http.ResponseWriter, r *http.Request) {
	getResponse(getNewsForCountryResponse, w, r.URL, true)
}
