package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type issInfo struct {
	Timestamp   int    `json:"timestamp"`
	Message     string `json:"message"`
	IssPosition struct {
		Longitude string `json:"longitude"`
		Latitude  string `json:"latitude"`
	} `json:"iss_position"`
}

func getISSPosition() ([]float64, error) {
	var i issInfo
	a := make([]float64, 2)

	response, err := http.Get("http://api.open-notify.org/iss-now.json")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode/100 != 2 {
		return nil, fmt.Errorf("bad response status: %s", response.Status)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}

	err = json.Unmarshal(responseData, &i)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal response body: %v", err)
	}

	if long, err := strconv.ParseFloat(i.IssPosition.Longitude, 64); err == nil {
		a[0] = long
	}
	if lat, err := strconv.ParseFloat(i.IssPosition.Latitude, 64); err == nil {
		a[1] = lat
	}

	return a, nil
}

func main() {
	pos, err := getISSPosition()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pos)
}
