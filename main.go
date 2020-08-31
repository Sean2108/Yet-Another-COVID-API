package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"yet-another-covid-map-api/casecount"
	"yet-another-covid-map-api/requests"
	"yet-another-covid-map-api/schedule"
)

var port string

func setupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/cases", requests.GetCaseCounts)
	http.HandleFunc("/news", requests.GetNewsForCountry)
}

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
}

func main() {
	// the John Hopkins data is updated at about 23:59 UTC everyday, so we will call update at 1am UTC
	schedule.CallFunctionDaily(casecount.UpdateCaseCounts, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	setupRoutes()
	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	wg.Wait()
}
