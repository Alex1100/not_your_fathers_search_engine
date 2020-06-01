package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	controllers "not_your_fathers_search_engine/api/controllers"

	gmux "github.com/gorilla/mux"
)

func setupServeMux() http.Handler {
	mux := gmux.NewRouter()

	mux.HandleFunc("/link", controllers.SearchLink).Methods("GET")

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

// InitializeApp initialize environment prior to app starting
func InitializeApp() {
	// loads values from config/.env.(current_env) into the system
	env := selectConfigFile()
	if err := godotenv.Load("config/" + env); err != nil {
		log.Print("No .env file found")
	}
}

// StartApp kick off the application once we load up main function
func StartApp() {
	projectID := os.Getenv("project_id")

	// Prints out projectId environment variable
	fmt.Println(projectID)

	mux := setupServeMux()
	fmt.Println("Listening on: ", 3010)

	log.Fatal(http.ListenAndServe(":3010", mux))
}
