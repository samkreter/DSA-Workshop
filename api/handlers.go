package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/samkreter/golang-utils/storage/influxdb"
)

const (
	database = "BitcoinPrice"
	username = "root"
	password = "root"
)

type QueryResult []SeriesEntry

type SeriesEntry struct {
	Price     json.Number `json:"price"`
	Timestamp string      `json:"timestamp"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://samkreter.github.io", http.StatusSeeOther)
}

func CheckAndConvertInputDate(date string) (string, bool) {
	if date == "" {
		return "", true
	}

	layout := "2006-01-02"
	timestamp, err := time.Parse(layout, date)
	if err != nil {
		return "", false
	}
	return timestamp.Format("2006-01-02 15:04:05"), true
}

func QueryInfluxDB(w http.ResponseWriter, r *http.Request) {
	urlParams := r.URL.Query()
	start, startOk := CheckAndConvertInputDate(urlParams.Get("start"))
	end, endOk := CheckAndConvertInputDate(urlParams.Get("end"))

	if start == "" {
		log.Println("Using Default Start: 2017-01-01")
		start = "2017-01-01"
	}

	if end == "" {
		log.Println("Using Default End: 2018-01-01")
		end = "2018-01-01"
	}

	if !startOk || !endOk {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusBadRequest, Text: "Start and end date must be in the form YYYY-MM-DD. Thanks and have a wonderful day."}); err != nil {
			log.Fatal(err)
		}
	}

	series := mux.Vars(r)["series"]

	var queryDb string

	//Make this safer by not taking user input directly
	switch series {
	case "bitcoin":
		queryDb = "Bitcoin"
	case "gold":
		queryDb = "Gold"
	case "silver":
		queryDb = "Silver"
	case "platinum":
		queryDb = "Platinum"
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusBadRequest, Text: "Could not find that databases."}); err != nil {
			log.Fatal(err)
		}
		return
	}

	c, err := influxDB.New(username, password, "influxdb")
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	var query string
	if start == "" || end == "" {
		query = fmt.Sprintf("Select price from %s where time > 2017-01-01 and time < 2018-01-01", queryDb)
	} else {
		query = fmt.Sprintf("Select price from %s where time > '%s' and time < '%s'", queryDb, start, end)
	}
	log.Println(query)
	metrics, err := c.ReadMetrics(database, query)
	if err != nil {
		log.Fatal(err)
	}

	var result QueryResult

	for _, entry := range metrics[0].Values {
		result = append(result, SeriesEntry{
			Price:     entry[1].(json.Number),
			Timestamp: entry[0].(string),
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Fatal(err)
	}

}

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
