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

func setupRoutes() {
	http.HandleFunc("/cases", requests.GetCaseCounts)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	// the John Hopkins data is updated at 11:50 or 11:55 pm GMT everyday, so we will call update at midnight utc (0 hour)
	schedule.CallFunctionDaily(casecount.UpdateCaseCounts, 0)
	wg := sync.WaitGroup{}
	wg.Add(1)
	setupRoutes()
	log.Printf("Server started at port%s", port)
	log.Fatalln(http.ListenAndServe(port, nil))
	wg.Wait()
}
