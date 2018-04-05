package main

import (
	"encoding/json"
	"fmt"
	"time"
	//"io"
	//"io/ioutil"
	"log"
	"net/http"

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

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos

*/
// func TodoCreate(w http.ResponseWriter, r *http.Request) {
// 	var todo Todo
// 	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err := r.Body.Close(); err != nil {
// 		panic(err)
// 	}
// 	if err := json.Unmarshal(body, &todo); err != nil {
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(422) // unprocessable entity
// 		if err := json.NewEncoder(w).Encode(err); err != nil {
// 			panic(err)
// 		}
// 	}

// 	t := RepoCreateTodo(todo)
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusCreated)
// 	if err := json.NewEncoder(w).Encode(t); err != nil {
// 		panic(err)
// 	}
// }
