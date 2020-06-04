package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

func insertData(position string) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO iss (position) VALUES ($1)", position)
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
	host := flag.String("host", "", "CrateDB hostname")
	port := flag.Int("port", 5432, "CrateDB postgresql port")
	flag.Parse()
	connStr := fmt.Sprintf("postgresql://crate@%s:%d/doc", *host, *port)

	var err error
	conn, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for {
		pos, err := getISSPosition()
		if err != nil {
			log.Fatal(err)
		} else {
			err = insertData(pos)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to insert data: %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Println("Sleeping for 5 seconds...")
		time.Sleep(time.Second * 5)
	}
}
