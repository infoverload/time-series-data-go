package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbpool *pgxpool.Pool

func insertData(position string) error {
	_, err := dbpool.Exec(context.Background(), "INSERT INTO iss (position) VALUES ($1)", position)
	return err
}

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

	s := fmt.Sprintf("POINT(%s %s)", i.IssPosition.Longitude, i.IssPosition.Latitude)
	return s, nil
}

func main() {
	var err error
	dbpool, err = pgxpool.Connect(context.Background(), "postgresql://crate@localhost:5433/doc")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	for {
		pos, err := getISSPosition()
		if err != nil {
			log.Fatal(err)
		}
		err = insertData(pos)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to insert data: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Sleeping for 10 seconds...")
		time.Sleep(time.Second * 10)
	}
}
