// Package main invokes the server package to start and run the REST API.
package main

import "GoMS1Assignment/restapi/server"

// main calls server.StartServer() to run the application.
func main() {
	server.StartServer()
}
