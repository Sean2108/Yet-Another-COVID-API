package requests

import (
	"net/url"
	"testing"
)

func TestParseUrlQuery(t *testing.T) {
	tables := []struct {
		rawurl   string
		query    string
		expected string
	}{
		{"http://localhost:8080/cases", "from", ""},
		{"http://localhost:8080/cases?from=&to=1/1/20", "from", ""},
		{"http://localhost:8080/cases?from=1/1/20&to=1/1/20", "from", "1/1/20"},
		{"http://localhost:8080/cases?from=&to=1/1/20", "to", "1/1/20"},
		{"http://localhost:8080/cases?from=&to=", "to", ""},
	}

	for _, table := range tables {
		url, _ := url.Parse(table.rawurl)
		result := parseURLQuery(url, table.query)
		if result != table.expected {
			t.Errorf("result of parseURLQuery was incorrect, got: %s, want: %s.", result, table.expected)
		}
	}

}
