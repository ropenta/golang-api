package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const (
	itemsPerPage = 20
)

type Station struct {
	ID             int    `json:"id"`
	StationName    string `json:"stationName"`
	AvailableDocks int    `json:"availableDocks"`
	TotalDocks     int    `json:"totalDocks"`
	StatusValue    string `json:"statusValue"`
	AvailableBikes int    `json:"availableBikes"`
	StAddress1     string `json:"stAddress1"`
}

type StationData struct {
	ExecutionTime   string    `json:"executionTime"`
	StationBeanList []Station `json:"stationBeanList"`
}

type DockableInfo struct {
	dockable bool
	message  string
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
	client := http.Client{}
	res, getErr := client.Do(stationReq)
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

func getStartAndEndIndices(startResults int, endResults int, pageInfo string) (start int, end int) {
	page, pageErr := strconv.Atoi(pageInfo)
	if pageInfo != "" {
		if pageErr != nil {
			fmt.Println("error converting string to int")
			return
		}
	}

	if page > 0 && (page-1)*itemsPerPage <= endResults {
		startResults += (page - 1) * itemsPerPage
		endResults = min(startResults+itemsPerPage, endResults)
	}
	return startResults, endResults
}

func getAllStations(w http.ResponseWriter, req *http.Request) {
	stations := getStations()

	// u, err := router.Get("YourHandler").URL("id", id, "key", key)

	startResults := 0
	endResults := len(stations)
	fmt.Println("GET params were:", req.URL.Query())
	fmt.Println(req.URL)
	pageInfo := req.URL.Query().Get("page")
	fmt.Println(reflect.TypeOf(pageInfo))

	startResults, endResults = getStartAndEndIndices(startResults, endResults, pageInfo)

	for i := startResults; i < endResults; i++ {
		response := fmt.Sprintf("%d: StationName: %s, TotalDocks: %d, StatusValue: %s, AvailableBikes: %d, Address: %s\n", i, stations[i].StationName, stations[i].TotalDocks, stations[i].StatusValue, stations[i].AvailableBikes, stations[i].StAddress1)
		fmt.Fprintf(w, response)
	}
	fmt.Println(stations[0].TotalDocks)
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
	startResults := 0
	endResults := len(stations)
	pageInfo := req.URL.Query().Get("page")
	startResults, endResults = getStartAndEndIndices(startResults, endResults, pageInfo)

	for i := startResults; i < endResults; i++ {
		response := fmt.Sprintf("%d: StationName: %s, TotalDocks: %d, StatusValue: %s, AvailableBikes: %d, Address: %s\n", i, stations[i].StationName, stations[i].TotalDocks, stations[i].StatusValue, stations[i].AvailableBikes, stations[i].StAddress1)
		fmt.Fprintf(w, response)
	}
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
	startResults := 0
	endResults := len(stations)
	pageInfo := req.URL.Query().Get("page")
	startResults, endResults = getStartAndEndIndices(startResults, endResults, pageInfo)

	for i := startResults; i < endResults; i++ {
		response := fmt.Sprintf("%d: StationName: %s, TotalDocks: %d, StatusValue: %s, AvailableBikes: %d, Address: %s\n", i, stations[i].StationName, stations[i].TotalDocks, stations[i].StatusValue, stations[i].AvailableBikes, stations[i].StAddress1)
		fmt.Fprintf(w, response)
	}
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

	for i := startResults; i < endResults; i++ {
		response := fmt.Sprintf("%d: StationName: %s, TotalDocks: %d, StatusValue: %s, AvailableBikes: %d, Address: %s\n", i, stations[i].StationName, stations[i].TotalDocks, stations[i].StatusValue, stations[i].AvailableBikes, stations[i].StAddress1)
		fmt.Fprintf(w, response)
	}
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
	message := fmt.Sprintf("There are %d available docks. You cannot return all %d of your bikes.", station.AvailableDocks, numBikesToReturn)
	if numBikesToReturn <= station.AvailableDocks {
		dockable = true
		message = fmt.Sprintf("There are %d available docks. You are able to return all %d of your bikes.", station.AvailableDocks, numBikesToReturn)
	}

	m := DockableInfo{dockable, message}
	fmt.Fprintf(w, fmt.Sprintf("%t\n", dockable))
	fmt.Fprintf(w, m.message)
}

func handleRequests() {
	router := mux.NewRouter()
	router.Methods("GET").Path("/stations").HandlerFunc(getAllStations)
	router.Methods("GET").Path("/stations").Queries("path", "{[0-9]*?}").HandlerFunc(getAllStations)
	router.Methods("GET").Path("/stations/in-service").HandlerFunc(getInServiceStations)
	router.Methods("GET").Path("/stations/in-service").Queries("path", "{[0-9]*?}").HandlerFunc(getInServiceStations)
	router.Methods("GET").Path("/stations/not-in-service").HandlerFunc(getNotInServiceStations)
	router.Methods("GET").Path("/stations/not-in-service").Queries("path", "{[0-9]*?}").HandlerFunc(getNotInServiceStations)
	router.Methods("GET").Path("/stations/{searchstring}").HandlerFunc(searchStations)
	router.Methods("GET").Path("/stations/{stationid}/{bikestoreturn}").HandlerFunc(returnBikes)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(router)

	if err := http.ListenAndServe(":4000", n); err != nil {
		log.Fatal(err)
	}
}

func main() {
	handleRequests()
}
