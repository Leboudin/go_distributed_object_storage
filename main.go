package main

import (
	"./api"

	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

// Initialize server storage root and objects folder,
// returns storage and error
func initStorage() (string, error) {
	// Get storage root path
	var storage = "/data"
	if os.Getenv("SERVER_ENV_STORAGE") != "" {
		storage = os.Getenv("SERVER_ENV_STORAGE")
	}
	log.Printf("Server storage root: %s", storage)

	err := os.Mkdir(storage, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Printf("Failed to create storage root: %s", storage)
		return storage, err
	}

	// Create objects folder
	err = os.Mkdir(storage+"/objects", os.ModePerm)
	if os.IsExist(err) {
		return storage, nil
	} else {
		return storage, err
	}
}

func main() {
	// Get server address config
	var addr = ":8030"
	if os.Getenv("SERVER_ENV_ADDR") != "" {
		addr = os.Getenv("SERVER_ENV_ADDR")
	}

	// First we need to initialize the storage, if failed, we have to exit
	storage, err := initStorage()
	if err != nil {
		log.Printf("Unable to initialize storage %s, error: %s", storage, err)
		log.Fatal("Now exiting...")
	}

	// Second, we initialize the API server
	apiSrv := api.NewServer(storage)

	// Routers
	router := httprouter.New()
	router.GET("/", apiSrv.Index)                  // Index, returns the server version, status
	router.GET("/objects/:name", apiSrv.GetObject) // RESTful API, get object by name
	router.PUT("/objects/:name", apiSrv.PutObject) // RESTful API, put object by name

	// Start the server at specified address
	log.Printf("Starting server on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
