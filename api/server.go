package api

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"log"
	"math/rand"
	uuid2 "github.com/satori/go.uuid"
	"sync"
)

type Status int64

const (
	PENDING  Status = iota
	RUNNING
	STOPPING
)

// API server struct holds the API server information
type Server struct {
	// API server version
	version int64

	// API server status
	status Status

	// Data provider serve details
	dp map[string]DataProvider

	// mutex on dp
	mutex sync.Mutex
}

type DataProvider struct {
	// Data provider server ID, a UUID string
	id string

	// Data provider server address
	addr string

	// Last pinged
	lastPing int64
}

// Get object name by storage and name
func (s *Server) getObjectName(name string) string {
	return s.storage + "/objects/" + name
}

// Create and return API server instance
func NewServer(dp []string) *Server {
	dps := map[string]DataProvider{}
	for i := range dp {
		provider := newDataProvider(dp[i])
		dps[provider.id] = *provider
	}

	return &Server{
		version: int64(1),
		status:  RUNNING,
		dp:      dps,
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

// RESTful API, get object by name
func (s *Server) GetObject(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	log.Printf("Getting object by name: %s", name)
	objName := s.getObjectName(name)
	GetObjectByName(objName, w)
}

// RESTful API, put object by name
// First we will choose a data server randomly, then we PUT file to the chosen server
func (s *Server) PutObject(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	objName := s.getObjectName(name)
	PutObjectByName(objName, w, r)
	s.selectDataProvider()
}

// Return a new DataProvider provided its address
func newDataProvider(addr string) *DataProvider {
	uuid := uuid2.Must(uuid2.NewV4())

	return &DataProvider{
		id:       uuid.String(),
		addr:     addr,
		lastPing: int64(0),
	}
}

// Select a DataProvider randomly for incoming PUT operation
func (s *Server) selectDataProvider() DataProvider {
	s.mutex.Lock()
	dps := make([]DataProvider, 0)
	for _, dp := range s.dp {
		dps = append(dps, dp)
	}
	s.mutex.Unlock()

	i := rand.Intn(len(dps))
	log.Printf("Selected data provider server %s, addr: %s", dps[i].id, dps[i].addr)
	return dps[i]
}
