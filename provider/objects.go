package provider

import (
	"net/http"
	"os"
	"log"
	"io"
)

// The real handler to get an object by object name
func GetObjectByName(name string, w http.ResponseWriter) {
	file, err := os.Open(name)
	if err != nil {
		log.Printf("Unable to open file %s, error: %s", name, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer file.Close()
	io.Copy(w, file)
}

// The real handler to put an object by object name
func PutObjectByName(name string, w http.ResponseWriter, r *http.Request) {
	file, err := os.Create(name)
	if err != nil {
		log.Printf("Unable to create file %s, error: %s", name, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()
	io.Copy(file, r.Body)
	log.Printf("Created object %s", name)
}
