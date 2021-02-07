package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/urfave/negroni"
)

// Station {station name, address, # bikes available, total # of docks}
type Station struct {
	StationName    string `json:"stationName"`
	AddressLine    string `json:"stAddress1"`
	TotalDocks     string `json:"totalDocks"`
	AvailableDocks string `json:"availableDocks"`
}

func homePage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the API home page!")
}

func getStations(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Here are the stations:")

}

func handleRequests() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/stations", getStations)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)

	log.Fatal(http.ListenAndServe(":4000", n))
}

func main() {
	handleRequests()
}
