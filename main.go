package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type issInfo struct {
	IssPosition struct {
		Longitude string `json:"longitude"`
		Latitude  string `json:"latitude"`
	} `json:"iss_position"`
}

func getISSPosition() (string, error) {
	var i issInfo

	response, err := http.Get("http://api.open-notify.org/iss-now.json")
	if err != nil {
		return "", fmt.Errorf("unable to retrieve request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode/100 != 2 {
		return "", fmt.Errorf("bad response status: %s", response.Status)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %v", err)
	}

	err = json.Unmarshal(responseData, &i)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal response body: %v", err)
	}

	s := fmt.Sprintf("(%s, %s)", i.IssPosition.Longitude, i.IssPosition.Latitude)
	return s, nil
}

func main() {
	pos, err := getISSPosition()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pos)
}
