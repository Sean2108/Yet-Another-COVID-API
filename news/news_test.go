package news

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"sort"
	"testing"
)

type mockClient struct{}

var mockJSONResponseFn func() (*http.Response, error)

type ByArticleTitle []Article

func (a ByArticleTitle) Len() int {
	return len(a)
}

func (a ByArticleTitle) Less(i, j int) bool {
	return a[i].Title < a[j].Title
}

func (a ByArticleTitle) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (m *mockClient) Get(url string) (*http.Response, error) {
	return mockJSONResponseFn()
}

func defaultJSONResponse() (*http.Response, error) {
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
	mockJSONResponseFn = defaultJSONResponse
	result, _ := GetNews("", "", "sg")
	expected := []Article{
		Article{"Google News", "headline", "desc", "testUrl", "imgUrl", "2020-04-02T01:01:22Z"},
		Article{"Google News2", "headline2", "desc2", "testUrl2", "imgUrl2", "2020-04-02T02:01:22Z"},
	}
	sort.Sort(ByArticleTitle(result))
	for i, item := range result {
		if item != expected[i] {
			t.Errorf("Result of GetNews is incorrect, got: %+v, want: %+v.", item, expected[i])
		}
	}
}

func TestGetNews_ReadJSONFailed(t *testing.T) {
	expected := "test failure"
	mockJSONResponseFn = func() (*http.Response, error) {
		return nil, errors.New(expected)
	}
	_, err := GetNews("", "", "sg")
	if err == nil {
		t.Error("GetNews should have thrown an error but it didn't.")
	}
	if err.Error() != "test failure" {
		t.Errorf("GetNews threw a different error than expected: got: %s, want: %s.", err.Error(), expected)
	}
}

func TestGetNews_MalformedJSON(t *testing.T) {
	mockJSONResponseFn = func() (*http.Response, error) {
		jsonStr := `{"status":"ok","totalResults":19,"articles":[]`
		r := ioutil.NopCloser(bytes.NewReader([]byte(jsonStr)))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	_, err := GetNews("", "", "sg")
	if err == nil {
		t.Error("GetNews should have thrown an error but it didn't.")
	}
}

func TestFormatResponse(t *testing.T) {
	input := []inputArticle{
		inputArticle{inputSource{"google-news", "Google News"}, "ST", "headline", "desc", "testUrl", "imgUrl", "2020-04-02T01:01:22Z", "testcontent"},
		inputArticle{inputSource{"google-news3", "Google News2"}, "ST2", "headline2", "desc2", "testUrl2", "imgUrl2", "2020-04-02T02:01:22Z", "testcontent2"},
		inputArticle{inputSource{"google-news2", "Google News2"}, "ST2", "headline2", "desc2", "testUrl2", "imgUrl2", "2020-04-02T02:01:22Z", "testcontent2"},
	}
	expected := []Article{
		Article{"Google News", "headline", "desc", "testUrl", "imgUrl", "2020-04-02T01:01:22Z"},
		Article{"Google News2", "headline2", "desc2", "testUrl2", "imgUrl2", "2020-04-02T02:01:22Z"},
	}
	result := formatResponse(input)
	sort.Sort(ByArticleTitle(result))
	for i, item := range result {
		if item != expected[i] {
			t.Errorf("Item in result of formatResponse is different from expected, got: %+v, want: %+v.", item, expected[i])
		}
	}
}
