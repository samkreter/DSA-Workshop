package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = ":8080"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/api/metrics/{series}", QueryInfluxDB).Methods("GET")
	log.Fatal(http.ListenAndServe(port, router))
}
