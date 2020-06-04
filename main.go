package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbpool *pgxpool.Pool

func createTable() error {
	_, err := dbpool.Exec(context.Background(), "CREATE TABLE iss (timestamp TIMESTAMP GENERATED ALWAYS AS CURRENT_TIMESTAMP, position GEO_POINT)")
	return err
}

func insertData(position []float64) error {
	_, err := dbpool.Exec(context.Background(), "INSERT INTO iss(position) VALUES($1)", position)
	return err
}

func listData() error {
	rows, _ := dbpool.Query(context.Background(), "SELECT * FROM iss")

	for rows.Next() {
		var timestamp string
		var position string
		err := rows.Scan(&timestamp, &position)
		if err != nil {
			return err
		}
		fmt.Printf("%s, %s\n", timestamp, position)
	}
	return rows.Err()
}

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
	var err error
	dbpool, err = pgxpool.Connect(context.Background(), "postgresql://crate@localhost:5433/doc")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	err = createTable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create table: %v\n", err)
		os.Exit(1)
	}

	pos, err := getISSPosition()
	if err != nil {
		log.Fatal(err)
	}
	err = insertData(pos)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data: %v\n", err)
		os.Exit(1)
	}

	err = listData()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to list data: %v\n", err)
		os.Exit(1)
	}
}
