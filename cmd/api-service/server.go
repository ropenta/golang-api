package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/urfave/negroni"
)

type Station struct {
	StationName    string `json:"stationName"`
	TotalDocks     int    `json:"totalDocks"`
	AvailableBikes int    `json:"availableBikes"`
	StAddress1     string `json:"stAddress1"`
}

type StationData struct {
	ExecutionTime   string    `json:"executionTime"`
	StationBeanList []Station `json:"stationBeanList"`
}

func homePage(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the API home page!")
}

func getStations(w http.ResponseWriter, req *http.Request) {
	urlEndpoint := "https://feeds.citibikenyc.com/stations/stations.json"
	req, err := http.NewRequest(http.MethodGet, urlEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	stationData := StationData{}
	jsonErr := json.Unmarshal(body, &stationData)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}
	stations := stationData.StationBeanList
	for i := 0; i < len(stations); i++ {
		response := fmt.Sprintf("StationName: %s, TotalDocks: %d, AvailableBikes: %d, Address: %s\n", stations[i].StationName, stations[i].TotalDocks, stations[i].AvailableBikes, stations[i].StAddress1)
		fmt.Fprintf(w, response)
	}
	fmt.Println(stations[0].TotalDocks)
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
