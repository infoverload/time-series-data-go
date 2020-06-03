package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

type positionISS struct {
	Timestamp   int    `json:"timestamp"`
	Message     string `json:"message"`
	IssPosition struct {
		Longitude string `json:"longitude"`
		Latitude  string `json:"latitude"`
	} `json:"iss_position"`
}

func pointWKT() (string, error) {

    var p positionISS

    response, err := http.Get("http://api.open-notify.org/iss-now.json")
    if err != nil {
        return "", fmt.Errorf("unable to retrieve request: %v", err)
    }
    defer response.Body.Close()

    responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return "", fmt.Errorf("unable to read response body: %v", err)
    }

    err = json.Unmarshal(responseData, &p)
    if err != nil {
        return "", fmt.Errorf("unable to unmarshal response body: %v", err)
    }

    return "POINT ( " + p.IssPosition.Longitude + " " + p.IssPosition.Latitude + " )", nil
}

func main() {
    pos, err := pointWKT()
    if err != nil {
      log.Fatal(err)
    }
    fmt.Println(pos)  
}