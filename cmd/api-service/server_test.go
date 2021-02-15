package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var allStationsJson = `{"executionTime":"2016-01-22 04:32:49 PM","stationBeanList":[{"id":72,"stationName":"W 52 St & 11 Ave","availableDocks":32,"totalDocks":39,"latitude":40.76727216,"longitude":-73.99392888,"statusValue":"In Service","statusKey":1,"availableBikes":7,"stAddress1":"W 52 St & 11 Ave","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2016-01-22 04:30:15 PM","landMark":""},{"id":423,"stationName":"W 54 St & 9 Ave","availableDocks":3,"totalDocks":3,"latitude":40.76584941,"longitude":-73.98690506,"statusValue":"Not In Service","statusKey":3,"availableBikes":0,"stAddress1":"W 54 St & 9 Ave","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2015-12-14 11:04:17 AM","landMark":""},{"id":79,"stationName":"Franklin St & W Broadway","availableDocks":0,"totalDocks":33,"latitude":40.71911552,"longitude":-74.00666661,"statusValue":"In Service","statusKey":1,"availableBikes":33,"stAddress1":"Franklin St & W Broadway","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2016-01-22 04:32:41 PM","landMark":""},{"id":82,"stationName":"St James Pl & Pearl St","availableDocks":27,"totalDocks":27,"latitude":40.71117416,"longitude":-74.00016545,"statusValue":"In Service","statusKey":1,"availableBikes":0,"stAddress1":"St James Pl & Pearl St","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2016-01-22 04:29:41 PM","landMark":""},{"id":83,"stationName":"Atlantic Ave & Fort Greene Pl","availableDocks":21,"totalDocks":62,"latitude":40.68382604,"longitude":-73.97632328,"statusValue":"In Service","statusKey":1,"availableBikes":40,"stAddress1":"Atlantic Ave & Fort Greene Pl","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2016-01-22 04:32:33 PM","landMark":""},{"id":116,"stationName":"W 17 St & 8 Ave","availableDocks":19,"totalDocks":39,"latitude":40.74177603,"longitude":-74.00149746,"statusValue":"In Service","statusKey":1,"availableBikes":19,"stAddress1":"W 17 St & 8 Ave","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2016-01-22 04:32:32 PM","landMark":""}]}`
var inServiceStationsJson = `{"executionTime":"2016-01-22 04:32:49 PM","stationBeanList":[{"id":72,"stationName":"W 52 St & 11 Ave","availableDocks":32,"totalDocks":39,"latitude":40.76727216,"longitude":-73.99392888,"statusValue":"In Service","statusKey":1,"availableBikes":7,"stAddress1":"W 52 St & 11 Ave","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2016-01-22 04:30:15 PM","landMark":""},{"id":79,"stationName":"Franklin St & W Broadway","availableDocks":0,"totalDocks":33,"latitude":40.71911552,"longitude":-74.00666661,"statusValue":"In Service","statusKey":1,"availableBikes":33,"stAddress1":"Franklin St & W Broadway","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2016-01-22 04:32:41 PM","landMark":""},{"id":82,"stationName":"St James Pl & Pearl St","availableDocks":27,"totalDocks":27,"latitude":40.71117416,"longitude":-74.00016545,"statusValue":"In Service","statusKey":1,"availableBikes":0,"stAddress1":"St James Pl & Pearl St","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2016-01-22 04:29:41 PM","landMark":""},{"id":83,"stationName":"Atlantic Ave & Fort Greene Pl","availableDocks":21,"totalDocks":62,"latitude":40.68382604,"longitude":-73.97632328,"statusValue":"In Service","statusKey":1,"availableBikes":40,"stAddress1":"Atlantic Ave & Fort Greene Pl","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2016-01-22 04:32:33 PM","landMark":""},{"id":116,"stationName":"W 17 St & 8 Ave","availableDocks":19,"totalDocks":39,"latitude":40.74177603,"longitude":-74.00149746,"statusValue":"In Service","statusKey":1,"availableBikes":19,"stAddress1":"W 17 St & 8 Ave","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2016-01-22 04:32:32 PM","landMark":""}]}`
var notInServiceStationsJson = `{"executionTime":"2016-01-22 04:32:49 PM","stationBeanList":[{"id":423,"stationName":"W 54 St & 9 Ave","availableDocks":3,"totalDocks":3,"latitude":40.76584941,"longitude":-73.98690506,"statusValue":"Not In Service","statusKey":3,"availableBikes":0,"stAddress1":"W 54 St & 9 Ave","stAddress2":"","city":"","postalCode":"","location":"","altitude":"","testStation":false,"lastCommunicationTime":"2015-12-14 11:04:17 AM","landMark":""}]}`

func init() {
	Client = &MockClient{}
}

type MockClient struct{}

var (
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func TestMin(t *testing.T) {
	a := min(1, 2)
	if a != 1 {
		t.Errorf("Wrong min, got %d, wanted %d", a, 2)
	}
	b := min(2, 1)
	if b != 1 {
		t.Errorf("Wrong min, got %d, wanted %d", b, 2)
	}
}

func TestGetStartAndEndIndices(t *testing.T) {
	startValue, endValue := getStartAndEndIndices(6, "1")
	if startValue != 0 && endValue != 6 {
		t.Errorf("Expected %d start value and %d end value, but received %d start and %d end", 0, 6, startValue, endValue)
	}
	startValue, endValue = getStartAndEndIndices(6, "invalid_page")
	if startValue != 0 && endValue != 6 {
		t.Errorf("Expected %d start value and %d end value, but received %d start and %d end", 0, 6, startValue, endValue)
	}
}

func TestGetStations(t *testing.T) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(allStationsJson)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	stations := getStations()
	if len(stations) != 6 {
		t.Errorf("Expected %d stations, but received %d stations", 6, len(stations))
	}
}

func TestGetAllStations(t *testing.T) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(allStationsJson)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	req, err := http.NewRequest("GET", "/stations", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllStations)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v but wanted %v", status, http.StatusOK)
	}

	expected := `[
    {
        "stationName": "W 52 St \u0026 11 Ave",
        "availableDocks": 32,
        "totalDocks": 39,
        "stAddress1": "W 52 St \u0026 11 Ave"
    },
    {
        "stationName": "W 54 St \u0026 9 Ave",
        "availableDocks": 3,
        "totalDocks": 3,
        "stAddress1": "W 54 St \u0026 9 Ave"
    },
    {
        "stationName": "Franklin St \u0026 W Broadway",
        "totalDocks": 33,
        "stAddress1": "Franklin St \u0026 W Broadway"
    },
    {
        "stationName": "St James Pl \u0026 Pearl St",
        "availableDocks": 27,
        "totalDocks": 27,
        "stAddress1": "St James Pl \u0026 Pearl St"
    },
    {
        "stationName": "Atlantic Ave \u0026 Fort Greene Pl",
        "availableDocks": 21,
        "totalDocks": 62,
        "stAddress1": "Atlantic Ave \u0026 Fort Greene Pl"
    },
    {
        "stationName": "W 17 St \u0026 8 Ave",
        "availableDocks": 19,
        "totalDocks": 39,
        "stAddress1": "W 17 St \u0026 8 Ave"
    }
]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetInServiceStations(t *testing.T) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(inServiceStationsJson)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	req, err := http.NewRequest("GET", "/stations/in-service", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getInServiceStations)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v but wanted %v", status, http.StatusOK)
	}

	expected := `[
    {
        "stationName": "W 52 St \u0026 11 Ave",
        "availableDocks": 32,
        "totalDocks": 39,
        "stAddress1": "W 52 St \u0026 11 Ave"
    },
    {
        "stationName": "Franklin St \u0026 W Broadway",
        "totalDocks": 33,
        "stAddress1": "Franklin St \u0026 W Broadway"
    },
    {
        "stationName": "St James Pl \u0026 Pearl St",
        "availableDocks": 27,
        "totalDocks": 27,
        "stAddress1": "St James Pl \u0026 Pearl St"
    },
    {
        "stationName": "Atlantic Ave \u0026 Fort Greene Pl",
        "availableDocks": 21,
        "totalDocks": 62,
        "stAddress1": "Atlantic Ave \u0026 Fort Greene Pl"
    },
    {
        "stationName": "W 17 St \u0026 8 Ave",
        "availableDocks": 19,
        "totalDocks": 39,
        "stAddress1": "W 17 St \u0026 8 Ave"
    }
]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetNotInServiceStations(t *testing.T) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(notInServiceStationsJson)))
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	req, err := http.NewRequest("GET", "/stations/not-in-service", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getNotInServiceStations)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v but wanted %v", status, http.StatusOK)
	}

	expected := `[
    {
        "stationName": "W 54 St \u0026 9 Ave",
        "availableDocks": 3,
        "totalDocks": 3,
        "stAddress1": "W 54 St \u0026 9 Ave"
    }
]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
