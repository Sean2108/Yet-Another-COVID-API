# Yet Another COVID API

[![Build Status](https://travis-ci.com/Sean2108/Yet-Another-COVID-Map-API.svg?branch=master)](https://travis-ci.com/Sean2108/Yet-Another-COVID-Map-API)[![codecov](https://codecov.io/gh/Sean2108/Yet-Another-COVID-Map-API/branch/master/graph/badge.svg)](https://codecov.io/gh/Sean2108/Yet-Another-COVID-Map-API)

API written in Go to provide number of COVID cases and deaths with date queries, and news stories related to COVID for countries.
Data is provided by John Hopkins CSSE (https://github.com/CSSEGISandData/COVID-19).

Deployed at https://yet-another-covid-api.herokuapp.com.

## Endpoints:
/cases:
- Call the endpoint with no query information (https://yet-another-covid-api.herokuapp.com/cases) to get the numbers of all confirmed cases and deaths for each state and country. 
- Call the endpoint with attributes 'from' and/or 'to' to get the numbers of all confirmed cases and deaths for each state and country between the from date and to date. For example https://yet-another-covid-api.herokuapp.com/cases?from=3/2/20&to=3/10/20.
- Call the endpoint with attribute 'aggregateCountries' set to true to aggregate the counts to the country level instead of the state level. For example, https://yet-another-covid-api.herokuapp.com/cases?aggregateCountries=true.
- Call the endpoint with country name in the field 'country' to extract the numbers of confirmed cases and deaths for all states in the country. For example, https://yet-another-covid-api.herokuapp.com/cases?country=Singapore.
- Call the endoint with attribute 'perDay' set to true to get a returned value without aggregated counts. Each state/country in the response will have a list of days with the cumulative number of confirmed cases and deaths up till that day. For example, https://yet-another-covid-api.herokuapp.com/cases?perDay=true.
- Call the endoint with attribute 'worldTotal' set to true to get the per day statistics for the world between the given from and to dates. For example, https://yet-another-covid-api.herokuapp.com/cases?from=3/2/20&to=3/10/20&worldTotal=true.

/news:
- Get news for country in the field to extract the latest coronavirus news for that country. Will use the News API (https://newsapi.org/) for obtaining this information. For example, https://yet-another-covid-api.herokuapp.com/news?country=Singapore
- Call the endpoint with attributes 'from' and/or 'to' to get the news between the from date and to date. For example https://yet-another-covid-api.herokuapp.com/news?from=3/2/20&to=3/10/20&country=us.

### Allowed date formats:
- MM/DD/YY
- MM/DD/YYYY
- YYYY/MM/DD
- YY/MM/DD

You can use either / or - as the date delimiters.

### Allowed country formats:
You can use the full name or the short 2 letter ISO 3166 Alpha-2 code to identify countries. For example, SG and Singapore are equivalent. This is case insensitive.
