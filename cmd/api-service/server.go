package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

const (
	itemsPerPage = 20
)

type Station struct {
	ID             int    `json:"id,omitempty"`
	StationName    string `json:"stationName"`
	AvailableDocks int    `json:"availableDocks,omitempty"`
	TotalDocks     int    `json:"totalDocks"`
	StatusValue    string `json:"statusValue,omitempty"`
	AvailableBikes int    `json:"availableBikes"`
	StAddress1     string `json:"stAddress1"`
}

type StationData struct {
	ExecutionTime   string    `json:"executionTime"`
	StationBeanList []Station `json:"stationBeanList"`
}

type DockableInfo struct {
	Dockable bool   `json:"dockable"`
	Message  string `json:"message"`
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getStations() []Station {
	urlEndpoint := "https://feeds.citibikenyc.com/stations/stations.json"
	stationReq, urlErr := http.NewRequest(http.MethodGet, urlEndpoint, nil)
	if urlErr != nil {
		log.Fatal(urlErr)
	}
	res, getErr := Client.Do(stationReq)
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
	return stations
}

func getStartAndEndIndices(numStations int, pageInfo string) (start int, end int) {
	startResults := 0
	endResults := numStations
	page, pageErr := strconv.Atoi(pageInfo)
	if pageInfo != "" {
		if pageErr != nil {
			fmt.Println("error converting string to int")
			return startResults, endResults
		}
	}

	if page > 0 && (page-1)*itemsPerPage <= endResults {
		startResults += (page - 1) * itemsPerPage
		endResults = min(startResults+itemsPerPage, endResults)
	}
	return startResults, endResults
}

func buildStationArry(stations []Station, startResults int, endResults int) []Station {
	var stationInfo []Station
	for i := startResults; i < endResults; i++ {
		station := stations[i]
		// remove these fields from view in JSON
		station.AvailableDocks = 0
		station.ID = 0
		station.StatusValue = ""
		stationInfo = append(stationInfo, station)
	}
	return stationInfo
}

func getAllStations(w http.ResponseWriter, req *http.Request) {
	stations := getStations()
	pageInfo := req.URL.Query().Get("page")
	startResults, endResults := getStartAndEndIndices(len(stations), pageInfo)

	var stationInfo []Station = buildStationArry(stations, startResults, endResults)
	stationMarshal, marshalErr := json.MarshalIndent(stationInfo, "", "    ")
	if marshalErr != nil {
		log.Fatal("Error marshaling struct to JSON")
	}
	fmt.Fprint(w, string(stationMarshal))
}

func getInServiceStations(w http.ResponseWriter, req *http.Request) {
	allStations := getStations()
	var inServiceStations []Station
	for _, v := range allStations {
		if v.StatusValue == "In Service" {
			inServiceStations = append(inServiceStations, v)
		}
	}

	stations := inServiceStations
	pageInfo := req.URL.Query().Get("page")
	startResults, endResults := getStartAndEndIndices(len(stations), pageInfo)

	var stationInfo []Station = buildStationArry(stations, startResults, endResults)
	stationMarshal, marshalErr := json.MarshalIndent(stationInfo, "", "    ")
	if marshalErr != nil {
		log.Fatal("Error marshaling struct to JSON")
	}
	fmt.Fprint(w, string(stationMarshal))
}

func getNotInServiceStations(w http.ResponseWriter, req *http.Request) {
	allStations := getStations()
	var notInServiceStations []Station
	for _, v := range allStations {
		if v.StatusValue == "Not In Service" {
			notInServiceStations = append(notInServiceStations, v)
		}
	}

	stations := notInServiceStations
	pageInfo := req.URL.Query().Get("page")
	startResults, endResults := getStartAndEndIndices(len(stations), pageInfo)

	var stationInfo []Station = buildStationArry(stations, startResults, endResults)
	stationMarshal, marshalErr := json.MarshalIndent(stationInfo, "", "    ")
	if marshalErr != nil {
		log.Fatal("Error marshaling struct to JSON")
	}
	fmt.Fprint(w, string(stationMarshal))
}

func searchStations(w http.ResponseWriter, req *http.Request) {
	allStations := getStations()
	var matchingStations []Station
	searchstring := strings.ToLower(mux.Vars(req)["searchstring"])

	for _, v := range allStations {
		if strings.Contains(strings.ToLower(v.StAddress1), searchstring) || strings.Contains(strings.ToLower(v.StationName), searchstring) {
			matchingStations = append(matchingStations, v)
		}
	}

	stations := matchingStations
	startResults := 0
	endResults := len(stations)

	var stationInfo []Station = buildStationArry(stations, startResults, endResults)
	stationMarshal, marshalErr := json.MarshalIndent(stationInfo, "", "    ")
	if marshalErr != nil {
		log.Fatal("Error marshaling struct to JSON")
	}
	fmt.Fprint(w, string(stationMarshal))
}

func returnBikes(w http.ResponseWriter, req *http.Request) {
	stationID := strings.ToLower(mux.Vars(req)["stationid"])
	numBikesToReturn, numError := strconv.Atoi(mux.Vars(req)["bikestoreturn"])
	if numError != nil {
		fmt.Println("error converting string to int")
		return
	}
	stations := getStations()

	station := Station{}
	for _, v := range stations {
		if strconv.Itoa(v.ID) == stationID {
			station = v
			break
		}
	}
	if strconv.Itoa(station.ID) == "" {
		log.Fatal("Station not found. Please enter a valid station id.")
	}
	dockable := false
	message := fmt.Sprintf("You cannot return all %d of your bikes. There are %d available docks.", numBikesToReturn, station.AvailableDocks)
	if station.StatusValue == "Not In Service" {
		message = fmt.Sprintf("Station %s with ID %d is Not In Service. Please choose an In Service station.", station.StationName, station.ID)
	} else if numBikesToReturn <= station.AvailableDocks {
		dockable = true
		message = fmt.Sprintf("You are able to return all %d of your bikes. There are %d available docks.", numBikesToReturn, station.AvailableDocks)
	}

	dockableMarshal, err := json.MarshalIndent(DockableInfo{Dockable: dockable, Message: message}, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling struct to JSON")
	}
	fmt.Fprint(w, string(dockableMarshal))
}
