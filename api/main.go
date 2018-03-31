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

	router := NewRouter()
	log.Println("Listing on port ", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}
