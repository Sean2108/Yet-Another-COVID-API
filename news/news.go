package news

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"yet-another-covid-map-api/utils"
)

var client utils.HTTPClient

// Article : the relevant information of each article that was retrieved
type Article struct {
	Source       string `json:"source"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
	PublishedAt  string `json:"publishedAt"`
}

type inputSource struct {
	ID   string
	Name string
}

type inputArticle struct {
	Source      inputSource
	Author      string
	Title       string
	Description string
	URL         string
	URLToImage  string
	PublishedAt string
	Content     string
}

type newsResponse struct {
	Status       string
	TotalResults int
	Articles     []inputArticle
}

// Please populate this field with your own News API key
var apiKey string

const (
	newsEnvironmentVar  string = "NEWS_API_KEY"
	newsAPIHeadlinesURL string = "https://newsapi.org/v2/top-headlines"
)

func init() {
	apiKey = os.Getenv(newsEnvironmentVar)
	if apiKey == "" {
		log.Fatal("News API key is not populated! Please add your apiKey to your " + newsEnvironmentVar + " environment variable.")
	}
	client = &http.Client{}
}

func readJSONFromURL(url string) (newsResponse, error) {
	log.Printf("calling News API at: %s\n", url)
	r, err := client.Get(url)

	response := newsResponse{}
	if err != nil {
		return response, err
	}
	defer r.Body.Close()

	decodeErr := json.NewDecoder(r.Body).Decode(&response)
	return response, decodeErr
}

func formSingleURLQuery(queryName string, value string) string {
	if value != "" {
		return fmt.Sprintf("&%s=%s", queryName, value)
	}
	return ""
}

func formURLQuery(from string, to string, country string) string {
	return fmt.Sprintf("%s?apiKey=%s&q=virus&language=en%s%s%s", newsAPIHeadlinesURL, apiKey,
		formSingleURLQuery("from", from), formSingleURLQuery("to", to), formSingleURLQuery("country", country))
}

func formatArticle(input inputArticle, ch chan Article, wg *sync.WaitGroup) {
	ch <- Article{input.Source.Name, input.Title, input.Description, input.URL, input.URLToImage, input.PublishedAt}
	wg.Done()
}

func formatResponse(input []inputArticle) []Article {
	numRows := len(input)
	ch := make(chan Article, numRows)
	wg := sync.WaitGroup{}
	for index := 0; index < numRows; index++ {
		wg.Add(1)
		go formatArticle(input[index], ch, &wg)
	}
	wg.Wait()
	close(ch)
	var result []Article
	set := make(map[string]bool)
	for item := range ch {
		if _, ok := set[item.Title]; !ok {
			set[item.Title] = true
			result = append(result, item)
		}
	}
	return result
}

// GetNews : get coronavirus related headlines for the country passed in the parameter and return them
func GetNews(from string, to string, country string) ([]Article, error) {
	urlQuery := formURLQuery(from, to, strings.ToLower(country))
	response, err := readJSONFromURL(urlQuery)
	return formatResponse(response.Articles), err
}
