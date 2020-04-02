package dateformat

import "time"

const (
	// CasesDateFormat : the date format used in the CSV file provided by John Hopkins CSSE
	CasesDateFormat = "1/2/06"
	// NewsDateFormat : the required date format used when making queries to the news API
	NewsDateFormat = "2006-01-02"
)

var dateFormats []string = []string{
	CasesDateFormat,
	NewsDateFormat,
	"1-2-06",
	"01/02/06",
	"01-02-06",
	"01/02/2006",
	"01-02-2006",
	"06/01/02",
	"06-1-2",
	"06/1/2",
	"2006-1-2",
	"2006/1/2",
	"06-01-02",
	"06/01/02",
}

// FormatDate : parse and format date into formatTo
func FormatDate(formatTo string, date string) (string, bool) {
	if date == "" {
		return "", true
	}
	for _, format := range dateFormats {
		if formattedDate, err := time.Parse(format, date); err == nil {
			return formattedDate.Format(formatTo), true
		}
	}
	return "", false
}
