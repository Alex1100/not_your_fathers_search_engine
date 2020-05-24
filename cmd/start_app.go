package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"


	gmux "github.com/gorilla/mux"
	controllers "not_your_fathers_search_engine/api/controllers"
)

func setupServeMux(rootController *controllers.CockRoachDataBase) http.Handler {
	mux := gmux.NewRouter()

	mux.HandleFunc("/link", rootController.SearchLink).Methods("GET")
	mux.HandleFunc("/link", rootController.UpsertLink).Methods("POST")

	mux.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/public")))
	http.Handle("/", mux)
	return mux
}

func selectConfigFile() string {
	// loads values from .env into the system
	env := ".env"
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	currentEnv := os.Getenv("env")

	if currentEnv == "development" {   
		env += ".development"
	} else if currentEnv == "production" {
		env += ".production"
	} else if currentEnv == "staging" {
		env += ".staging"
	}

	return env
}

func InitializeApp() {
	// loads values from config/.env.(current_env) into the system
	env := selectConfigFile()
	if err := godotenv.Load("config/" + env); err != nil {
		log.Print("No .env file found")
	}
}

func StartApp() {
	projectId := os.Getenv("project_id")

	// Prints out username environment variable
	fmt.Println(projectId)

	rootController := controllers.ExposeDB()
	defer rootController.DB.DB.Close()

	mux := setupServeMux(rootController)
	fmt.Println("Listening on: ", 3010)

	log.Fatal(http.ListenAndServe(":3010", mux))
}