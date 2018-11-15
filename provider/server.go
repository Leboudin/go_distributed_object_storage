package provider

import (
	"os"
	"log"
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

func NewServer(addr string, storage string) *DataProviderServer {
	err := initStorage(storage)
	if err != nil {
		log.Printf("Unable to initialize storage %s, error: %s", storage, err)
		log.Fatal("Data provider server " + addr + "exiting...")
	}

	return &DataProviderServer{
		version: int64(1),
		addr:    addr,
		storage: storage,
	}
}
