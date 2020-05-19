package main

import (
	"fmt"

	gmux "github.com/gorilla/mux"
	controllers "not_your_fathers_search_engine/controllers"

	"log"
	"net/http"
)

// WORK IN PROGRESS

func setupServeMux(rootController *controllers.CockRoachDataBase) http.Handler {
	mux := gmux.NewRouter()

	mux.HandleFunc("/link", rootController.SearchLink).Methods("GET")
	mux.HandleFunc("/link", rootController.UpsertLink).Methods("POST")

	mux.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/public")))
	http.Handle("/", mux)
	return mux
}


func main() {
	rootController := controllers.ExposeDB()
	defer rootController.DB.DB.Close()

	mux := setupServeMux(rootController)
	defer
	fmt.Println("Listening on: ", 3010)

	log.Fatal(http.ListenAndServe(":3010", mux))
}