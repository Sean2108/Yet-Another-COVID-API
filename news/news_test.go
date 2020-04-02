package news

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type mockClient struct{}

func (m *mockClient) Get(url string) (*http.Response, error) {
	jsonStr := `{"status":"ok","totalResults":19,"articles":[
		{"source":{"id":"google-news","name":"Google News"},"author":"ST","title":"headline","description":"desc","url":"testUrl","urlToImage":"imgUrl","publishedAt":"2020-04-02T01:01:22Z","content":"testcontent"}
		{"source":{"id":"google-news2","name":"Google News2"},"author":"ST2","title":"headline2","description":"desc2","url":"testUrl2","urlToImage":"imgUrl2","publishedAt":"2020-04-02T02:01:22Z","content":"testcontent2"}
		]}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonStr)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

func TestFormURLQuery(t *testing.T) {
	apiKey = "testkey"
	tables := []struct {
		from     string
		to       string
		country  string
		expected string
	}{
		{"2020-01-02", "2020-01-03", "", "https://newsapi.org/v2/top-headlines?apiKey=testkey&q=virus&language=en&from=2020-01-02&to=2020-01-03"},
		{"", "", "sg", "https://newsapi.org/v2/top-headlines?apiKey=testkey&q=virus&language=en&country=sg"},
		{"2020-01-02", "", "us", "https://newsapi.org/v2/top-headlines?apiKey=testkey&q=virus&language=en&from=2020-01-02&country=us"},
	}

	for _, table := range tables {
		result := formURLQuery(table.from, table.to, table.country)
		if result != table.expected {
			t.Errorf("Result of formUrlQuery was incorrect, got: %s, want: %s.", result, table.expected)
		}
	}
}

func TestGetNews(t *testing.T) {
	client = &mockClient{}
	result, _ := GetNews("", "", "sg")
	expected := []Article{
		Article{"Google News", "headline", "desc", "testUrl", "imgUrl", "2020-04-02T01:01:22Z"},
		Article{"Google News2", "headline2", "desc2", "testUrl2", "imgUrl2", "2020-04-02T02:01:22Z"},
	}
	for i, item := range result {
		if item != expected[i] {
			t.Errorf("Result of GetNews is incorrect, got: %+v, want: %+v.", item, expected[i])
		}
	}
}
