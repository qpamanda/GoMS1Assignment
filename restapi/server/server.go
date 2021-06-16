/*
Package server initialises the handler functions for the REST API
and implements its functions for CRUD operations.
It also initialises the database for interfacing with the database package.
It is separated into 2 .go files to segregate the functionalities of the application.

	server.go: Initialises the handler functions and database, then starts the REST API to run
	on the designated port.

	handler.go: Implements the functions for CRUD operations as called by the client.
	Interface with the database package for database operations.
*/
package server

import (
	"GoMS1Assignment/restapi/database"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// StartServer initialises the database and handler functions then
// listens on the designated port to start the REST API running.
func StartServer() {
	// Initialise the database
	initDB()

	router := mux.NewRouter()

	// Initialise the handlers
	initaliseHandlers(router)

	// Set the listen port
	fmt.Println("Listening at port 5000")
	err := http.ListenAndServe(":5000", router)
	if err != nil {
		log.Fatal("FATAL: ListenAndServe - ", err)
	}

	defer database.DB.Close()
}

// initaliseHandlers initialises the handlers for the REST API.
func initaliseHandlers(router *mux.Router) {
	router.HandleFunc("/api/v1/", home)
	router.HandleFunc("/api/v1/courses", allcourses)
	router.HandleFunc("/api/v1/courses/{courseid}", course).Methods("GET", "PUT", "POST", "DELETE")
}

// initDB initialises the database
func initDB() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(">> Panic:", err)
		}
	}()

	config := getDBConfig()

	// Get database connection string
	connectionString := database.GetConnectionString(config)

	// Connect to database
	err := database.Connect(connectionString)
	if err != nil {
		panic("error connecting to database")
	}

	// Test connection to database
	err = database.DB.Ping()
	if err != nil {
		panic("error pinging to database")
	} else {
		fmt.Println("Ping to database success")
	}
}

// getDBConfig retrieves the database configurations and returns a struct.
func getDBConfig() database.Config {
	// Load setup.env file from same directory
	err := godotenv.Load("setup.env")
	if err != nil {
		log.Fatal("FATAL: Error loading .env file")
	}

	// Get env variables for database configuration
	serverName := os.Getenv("SERVER_NAME")
	dbName := os.Getenv("DB_NAME")
	dbUserName := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	config :=
		database.Config{
			ServerName: serverName,
			User:       dbUserName,
			Password:   dbPassword,
			DB:         dbName,
		}

	return config
}
