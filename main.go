package main

import (
	"./api"

	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func main() {
	// Get server address config
	var addr = ":8030"
	if os.Getenv("SERVER_ENV_ADDR") != "" {
		addr = os.Getenv("SERVER_ENV_ADDR")
	}

	// First, we initialize the API server with data providers
	apiSrv := api.NewServer([]string{"localhost:8030", "localhost:8031"})

	// Routers
	router := httprouter.New()
	router.GET("/", apiSrv.Index)                  // Index, returns the server version, status
	router.GET("/objects/:name", apiSrv.GetObject) // RESTful API, get object by name
	router.PUT("/objects/:name", apiSrv.PutObject) // RESTful API, put object by name

	// Start the server at specified address
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
