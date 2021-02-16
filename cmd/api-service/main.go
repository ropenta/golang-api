package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func handleRequests() {
	router := mux.NewRouter()
	router.Methods("GET").Path("/stations").HandlerFunc(getAllStations)
	router.Methods("GET").Path("/stations/in-service").HandlerFunc(getInServiceStations)
	router.Methods("GET").Path("/stations/not-in-service").HandlerFunc(getNotInServiceStations)
	router.Methods("GET").Path("/stations/{searchstring}").HandlerFunc(searchStations)
	router.Methods("GET").Path("/stations/{stationid}/{bikestoreturn}").HandlerFunc(returnBikes)

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(router)

	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Starting server with port 4000")
	if err := http.ListenAndServe(":4000", n); err != nil {
		log.Fatal("Error starting server", err)
	}
}

func main() {
	handleRequests()
}
