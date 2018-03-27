package main

import (
	"log"
	"net/http"
)

const port = ":8080"

func main() {

	router := NewRouter()
	log.Println("Listing on port ", port)
	log.Fatal(http.ListenAndServe(port, router))
}
