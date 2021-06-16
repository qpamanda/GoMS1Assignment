/*
Package client initialises the handler functions for the client web pages
and implements its functions for CRUD operations.
It is separated into 3 .go files to segregate the functionalities of the application.

	client.go: Initialises the templates and handler functions, then starts the client to run
	on the designated port.

	handler.go: Implements the handler functions for displaying the web pages of the client
	to perform the CRUD operations.

	crud.go: Implements the functions for CRUD operations and invokes the REST API.
*/
package client

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	tpl *template.Template
)

// InitClient initialises the templates for displaying the web pages of the client
func InitClient() {
	// Parse templates
	tpl = template.Must(template.ParseGlob("templates/*"))
}

// StartClient initialises the handler functions then
// listens on the designated port to start the client running.
func StartClient() {
	router := mux.NewRouter()

	// Initialise the handlers
	initaliseHandlers(router)

	// Set the listen port
	fmt.Println("Listening at port 5221")
	err := http.ListenAndServeTLS(":5221", "certs//cert.pem", "certs//key.pem", router)
	if err != nil {
		log.Fatal("FATAL: ListenAndServeTLS - ", err)
	}
}

// initaliseHandlers initialises the handlers for the client.
func initaliseHandlers(router *mux.Router) {
	router.HandleFunc("/", index)
	router.HandleFunc("/addcourse", addcourse)
	router.HandleFunc("/updcourse", updcourse)
	router.HandleFunc("/delcourse", delcourse)
	router.Handle("/favicon.ico", http.NotFoundHandler())
}
