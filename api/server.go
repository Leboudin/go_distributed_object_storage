package api

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"log"
	"math/rand"
	uuid2 "github.com/satori/go.uuid"
	"sync"
	"errors"

	"../streams"
	"io"
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

// DataProvider stores DataProviderServer info inside Server instance,
// our data will be saved to DataProviderServer
type DataProvider struct {
	// Data provider server ID, a UUID string
	id string

	// Data provider server address
	addr string

	// Last pinged
	lastPing int64
}

// Create and return API server instance
func NewServer(dp []string) *Server {
	dps := map[string]DataProvider{}
	for i := range dp {
		if dp[i] != "" {
			provider := newDataProvider(dp[i])
			dps[provider.id] = *provider
		}
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
func (s *Server) selectDataProvider() (DataProvider, error) {
	s.mutex.Lock()
	dps := make([]DataProvider, 0)
	for _, dp := range s.dp {
		dps = append(dps, dp)
	}
	s.mutex.Unlock()

	if len(dps) == 0 {
		return DataProvider{}, errors.New("no data server available")
	}

	i := rand.Intn(len(dps))
	//log.Printf("Selected data provider server %s, addr: %s", dps[i].id, dps[i].addr)
	return dps[i], nil
}

// Get object from data provider server
func (s *Server) GetObject(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	dataSrv, err := s.selectDataProvider()
	if err != nil {
		log.Printf("Unable to select data server, error: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	objNameWithAddr := dataSrv.addr + "/objects/" + name
	getStream, err := streams.NewGetStream(objNameWithAddr)
	if err != nil {
		log.Printf("Failed to get object %s, error: %s", objNameWithAddr, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	io.Copy(w, getStream)
	log.Printf("Successfully get object %s from %s (%s)", name, dataSrv.id, dataSrv.addr)

}

// Put object to data provider server
func (s *Server) PutObject(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := p.ByName("name")
	dataSrv, err := s.selectDataProvider()
	if err != nil {
		log.Printf("Unable to select data server, error: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	objNameWithAddr := dataSrv.addr + "/objects/" + name
	putStream := streams.NewPutStream(objNameWithAddr)

	io.Copy(putStream, r.Body)
	err = putStream.Close()

	if err != nil {
		log.Printf("Failed to put object %s, error: %s", objNameWithAddr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully put object %s to data server %s (%s)", name, dataSrv.id, dataSrv.addr)
}
