package api

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"log"
)

type Status int64

const (
	PENDING  Status = iota
	RUNNING
	STOPPING
)

// API server struct holds the API server information
type Server struct {
	// Storage root path
	storage string

	// API server version
	version int64

	// API server status
	status Status
}

// Create and return API server instance
func NewServer(storage string) *Server {
	return &Server{
		storage: storage,
		version: int64(1),
		status:  RUNNING,
	}
}

// Serves "/" index page, returns API server info
func (s *Server) Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	info := map[string]interface{}{
		"version": s.version,
		"status":  s.status,
	}

	resp, _ := json.Marshal(info)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// Get object name by storage and name
func (s *Server) getObjectName(name string) string {
	return s.storage + "/objects/" + name
}

// RESTful API, get object by name
func (s *Server) GetObject(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	log.Printf("Getting object by name: %s", name)
	objName := s.getObjectName(name)
	GetObjectByName(objName, w)
}

// RESTful API, put object by name
func (s *Server) PutObject(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	objName := s.getObjectName(name)
	PutObjectByName(objName, w, r)
}
