package provider

import (
	"os"
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
)

type DataProviderServer struct {
	// Provider server version
	version int64

	// Provider server address
	addr string

	// Storage root path
	storage string
}

// Initialize server storage root and objects folder,
// returns storage and error
func initStorage(storage string) error {
	log.Printf("Data provider server storage root: %s", storage)

	err := os.Mkdir(storage, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Printf("Failed to create storage root: %s", storage)
		return err
	}

	// Create objects folder
	err = os.Mkdir(storage+"/objects", os.ModePerm)
	if os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

// Initialize storage and return DataProviderServer instance
func NewServer(addr string, storage string) *DataProviderServer {
	err := initStorage(storage)
	if err != nil {
		log.Printf("Unable to initialize storage %s, error: %s", storage, err)
		log.Fatal("Data provider server " + addr + " exiting...")
	}

	return &DataProviderServer{
		version: int64(1),
		addr:    addr,
		storage: storage,
	}
}

// Get object name by storage and name
func (s *DataProviderServer) getObjectName(name string) string {
	return s.storage + "/objects/" + name
}

// RESTful API, get object by name
func (s *DataProviderServer) GetObject(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	log.Printf("Getting object by name: %s", name)
	objName := s.getObjectName(name)
	GetObjectByName(objName, w)
}

// RESTful API, put object by name
// First we will choose a data server randomly, then we PUT file to the chosen server
func (s *DataProviderServer) PutObject(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	objName := s.getObjectName(name)
	PutObjectByName(objName, w, r)
}
