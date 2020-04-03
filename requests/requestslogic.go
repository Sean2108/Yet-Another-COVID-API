package requests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"yet-another-covid-map-api/casecount"
	"yet-another-covid-map-api/dateformat"
	"yet-another-covid-map-api/news"
)

type writer interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

func parseURL(URL *url.URL, getAbbreviation bool, dateFormat string) (string, string, string, bool, bool) {
	from := parseURLQuery(URL, "from")
	to := parseURLQuery(URL, "to")
	country := parseURLQuery(URL, "country")

	from, ok := dateformat.FormatDate(dateFormat, from)
	if !ok {
		return "", "", "", false, false
	}
	to, ok = dateformat.FormatDate(dateFormat, to)
	if !ok {
		return "", "", "", false, false
	}

	var countryLookupFuncToCall func(string) (string, bool)
	if getAbbreviation {
		countryLookupFuncToCall = getAbbreviationFromCountry
	} else {
		countryLookupFuncToCall = getCountryFromAbbreviation
	}
	if countryFromAbbr, ok := countryLookupFuncToCall(country); ok {
		country = countryFromAbbr
	}
	aggregateCountries := parseURLQuery(URL, "aggregateCountries") == "true"

	return from, to, country, aggregateCountries, true
}

func parseURLQuery(URL *url.URL, key string) string {
	keys, ok := URL.Query()[key]
	var result string
	if ok && len(keys) > 0 {
		result = keys[0]
	}
	return result
}

func getCaseCountsResponse(from string, to string, country string, aggregateCountries bool) ([]byte, error, error) {
	if aggregateCountries {
		caseCounts, caseCountsErr := casecount.GetCountryCaseCounts(from, to, country)
		if caseCountsErr != nil {
			return nil, nil, caseCountsErr
		}
		response, err := json.Marshal(caseCounts)
		return response, err, nil
	}
	caseCounts, caseCountsErr := casecount.GetCaseCounts(from, to, country)
	if caseCountsErr != nil {
		return nil, nil, caseCountsErr
	}
	response, err := json.Marshal(caseCounts)
	return response, err, nil
}

func getNewsForCountryResponse(from string, to string, country string, _ bool) ([]byte, error, error) {
	articles, newsErr := news.GetNews(from, to, country)
	response, err := json.Marshal(articles)
	return response, err, newsErr
}

func getResponse(getDataFn func(from string, to string, country string, aggregateCountries bool) ([]byte, error, error), w writer, URL *url.URL) {
	log.Println(URL.String())
	from, to, country, aggregateCountries, ok := parseURL(URL, false, dateformat.CasesDateFormat)
	if !ok {
		http.Error(w, "Date format is not recognised, please use either YYYY-MM-DD, YYYY/MM/DD, MM-DD-YY or MM/DD/YY", http.StatusBadRequest)
		return
	}
	response, jsonErr, internalErr := getDataFn(from, to, country, aggregateCountries)
	if internalErr != nil {
		http.Error(w, internalErr.Error(), http.StatusBadRequest)
		return
	}
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func writeLandingPageHTML(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(
		`<head>
		<title>
			Yet Another COVID Map API
		</title>
		<h1>
			Yet Another COVID Map API
		</h1>
	</head>
	<body>
		<h2>Endpoints</h2>
		<h3>/cases</h3>
		<ul>
			<li>
				Call the endpoint with no query information (<a href=https://yet-another-covid-api.herokuapp.com/cases>/cases</a>) to get the numbers of all confirmed cases and deaths for each state and country.
			</li>
			<li>
				Call the endpoint with attributes 'from' and/or 'to' to get the numbers of all confirmed cases and deaths for each state and country between the from date and to date. These dates should be in the format YYYY-MM-DD or M/D/YY, for example <a href=https://yet-another-covid-api.herokuapp.com/cases?from=3/2/20&to=3/10/20>/cases?from=3/2/20&to=3/10/20</a>.
			</li>
			<li>
				Call the endpoint with attribute 'aggregateCountries' set to true to aggregate the counts to the country level instead of the state level. For example, <a href=https://yet-another-covid-api.herokuapp.com/cases?aggregateCountries=true>/cases?aggregateCountries=true</a>.
			</li>
			<li>
				Call the endpoint with country name in the field 'country' to extract the numbers of confirmed cases and deaths for all states in the country. For example, <a href=https://yet-another-covid-api.herokuapp.com/cases?country=Singapore>/cases?country=Singapore</a>.
			</li>
		</ul>
		<h3>/news</h3>
		<ul>
			<li>
				Get news for country in the field to extract the latest coronavirus news for that country. Will use the News API (<a href=https://newsapi.org>https://newsapi.org</a>) for obtaining this information. For example, <a href=https://yet-another-covid-api.herokuapp.com/news?country=Singapore>/news?country=Singapore</a>.
			</li>
			<li>
				Call the endpoint with attributes 'from' and/or 'to' to get the news between the from date and to date. For example <a href=https://yet-another-covid-api.herokuapp.com/news?from=3/2/20&to=3/10/20&country=us>/news?from=3/2/20&to=3/10/20&country=us</a>.
			</li>
		</ul>
		<h4>Allowed date formats:</h4>
		<ul>		
			<li>MM/DD/YY</li>
			<li>MM/DD/YYYY</li>
			<li>YYYY/MM/DD</li>
			<li>YY/MM/DD</li>
		</ul>
		You can use either / or - as the date delimiters.
		<h4>Allowed country formats:</h4>
		You can use the full name or the short 2 letter ISO 3166 Alpha-2 code to identify countries.
	</body>`))
}
