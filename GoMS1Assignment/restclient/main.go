// Package main invokes the client package to start and run the client.
package main

import (
	"GoMS1Assignment/restclient/client"
)

// init calls client.InitClient to initialise the variables. This will only be called once
// in the duration of the application.
func init() {
	client.InitClient()
}

// main calls client.StartClient() to run the application.
func main() {
	client.StartClient()
}
