package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"

	"./api"
	"./provider"
	"./util"
)

func main() {
	addr := flag.String("address", ":8030", "The server will listen on this address")
	storage := flag.String("storage", "/data", "The storage path will be used to store files")
	dps := flag.String("dps", "",
		"The comma seperated ip address of data provider servers, e.g. \"localhost:8030,localhost:8031\"")
	flag.Parse()

	switch flag.Arg(0) {
	case "dataserver":
		startDataServer(*addr, *storage)
	case "server":
		startAPIServer(*addr, *dps)
	default:
		startAPIServer(*addr, *dps)
	}
}

func startAPIServer(addr string, dps string) {
	log.Printf("Starting API server on %s", addr)
	dpList := util.ProcessIP(dps)
	log.Printf("Data provider servers: %s", dpList)

	// We initialize the API server with data providers
	apiSrv := api.NewServer(dpList)

	// Routers
	router := httprouter.New()
	router.GET("/", apiSrv.Index)                  // Index, returns the server version, status
	router.GET("/objects/:name", apiSrv.GetObject) // RESTful API, get object by name
	router.PUT("/objects/:name", apiSrv.PutObject) // RESTful API, put object by name

	// Start serving
	log.Fatal(http.ListenAndServe(addr, router))
}

func startDataServer(addr string, storage string) {
	log.Printf("Starting data provider server on %s, storage root: %s", addr, storage)

	// We initialize the data server with addr and storage
	dataSrv := provider.NewServer(addr, storage)

	// Listen to the object location query queue
	go func() {
		dataSrv.ListenToObjectLocateQueue()
	}()

	// Routers
	router := httprouter.New()
	router.GET("/objects/:name", dataSrv.GetObject) // RESTful API, get object by name
	router.PUT("/objects/:name", dataSrv.PutObject) // RESTful API, put object by name

	// Start serving
	log.Fatal(http.ListenAndServe(addr, router))
}
