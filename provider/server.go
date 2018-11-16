package provider

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"

	"../sqs"
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

// Listens to object location query queue, consume messages from API server,
// the message looks like:
// 		{"name": "someobject", "uid": <UUID>}
// the data provider server will try to find the object in its storage path,
// if found, data provider server will send a message to "godos-test-located" queue,
//		{"uid": <UUID>, "addr": "localhost:8031", "name": "someobject"}
// 		(the uid is used to determine which request by API server)
// then deletes the message from object location query queue, to prevent it being consumed again
// if not found, ignore it
func (s *DataProviderServer) ListenToObjectLocateQueue() {
	locateSQS := sqs.NewSQS()
	defer locateSQS.Close()

	c := locateSQS.Consume(locateSQS.Url)
	for r := range c {
		log.Printf("Consume message %s", *r.Body)
		req := make(map[string]string)
		err := json.Unmarshal([]byte(*r.Body), &req)
		if err != nil {
			log.Printf("Failed to unmarshal string to JSON, error: %s", err)
			continue
		}

		log.Printf("Trying to located object %s, request UID: %s", req["name"], req["uid"])
		if s.isObjectExists(req["name"]) {
			log.Printf("Object %s found", req["name"])
			// delete the message
			go func() {
				err := locateSQS.DeleteMessage(r, locateSQS.Url)
				if err != nil {
					log.Printf("Failed to delete message %s, error: %s", *r.ReceiptHandle, err)
				}
			}()

			// send reply to godos-test-located queue
			go func() {
				msg := map[string]string{
					"name": req["name"],
					"uid":  req["uid"],
					"addr": s.addr,
				}

				_, err := locateSQS.SendMessage(msg, locateSQS.ReplyUrl)
				if err != nil {
					log.Printf("Failed to send reply message, error: %s", err)
				}
			}()
		} else {
			log.Printf("Object %s not found", req["name"])
		}
	}
}

// Determines if object exists
func (s *DataProviderServer) isObjectExists(name string) bool {
	_, err := os.Stat(s.storage + "/objects/" + name)
	return !os.IsNotExist(err)
}
