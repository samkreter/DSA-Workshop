package main

import (
	"encoding/json"
	"fmt"
	//"io"
	//"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samkreter/DSA-Workshop/storage/influxdb"
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
	//t := r.URL.Query()["test"]
	vars := mux.Vars(r)

	t := vars["test"]

	if t == "" {
		fmt.Fprint(w, "empty")
	} else {
		fmt.Fprint(w, t)
	}

	log.Println("done")
}

func QueryInfluxDB(w http.ResponseWriter, r *http.Request) {
	// query := r.URL.Query()["query"]

	// if query == nil {
	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusBadRequest, Text: "Must proide a query for the database"}); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

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

	c, err := influxDB.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	query := fmt.Sprintf("Select price from %s", queryDb)

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
