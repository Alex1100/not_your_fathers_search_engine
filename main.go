package main

import (
	"fmt"

	gmux "github.com/gorilla/mux"
	controllers "not_your_fathers_search_engine/controllers"

	"log"
	"net/http"
)

// WORK IN PROGRESS

func setupServeMux() http.Handler {
	apiController := controllers.NewApiController()
	mux := gmux.NewRouter()

	mux.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/public")))
	http.Handle("/", mux)
	mux.HandleFunc("/link", apiController.LinkResource.SearchLink).Methods("GET")
	mux.HandleFunc("/link", apiController.LinkResource.UpsertLink).Methods("POST")

	return mux
}


func main() {
	mux := setupServeMux()
	
	fmt.Println("Listening on: ", 3010)

	log.Fatal(http.ListenAndServe(":3010", mux))
}