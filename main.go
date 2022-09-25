package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var applicationStartTime time.Time

type KeyValEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KVList []KeyValEntry

func getUptime() time.Duration {
	return time.Since(applicationStartTime).Round(time.Second)
}

func getKeyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	keyRes, err := Get(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(*keyRes + "\n"))
	if err != nil {
		log.Printf("getKeyHandlerFunc - Error %s", err)
	}
}

func getAllKeyHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	contents := GetAll()
	kvlist, err := json.Marshal(contents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Error"))
		if err != nil {
			log.Printf("getAllKeyHandlerFunc - Error %s", err)
		}
	}
	_, err = w.Write(kvlist)
	if err != nil {
		log.Printf("getAllKeyHandlerFunc - Error %s", err)
	}
}

func addKeyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	var kvEntry KeyValEntry

	switch r.Method {
	// add new key
	case http.MethodPost:
		// if exists, error thrown
		_ = json.NewDecoder(r.Body).Decode(&kvEntry)

		err := Put(kvEntry.Key, kvEntry.Value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		jsonKv, err := json.Marshal(kvEntry)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Error"))
			if err != nil {
				log.Printf("addKeyHandlerFunc - Error %s", err)
			}
		}
		_, err = w.Write(jsonKv)
		if err != nil {
			log.Printf("addKeyHandlerFunc - Error %s", err)
		}

	case http.MethodPut:
		// update the key; if exists, put isn't valid
		_ = json.NewDecoder(r.Body).Decode(&kvEntry)

		err := Update(kvEntry.Key, kvEntry.Value)
		if err != nil {
			// key does not exist; cannot update
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		jsonKv, err := json.Marshal(kvEntry)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("Error"))
			if err != nil {
				log.Printf("Error: %s", err)
			}
		}
		_, err = w.Write(jsonKv)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
}

func baseHandlerFunc(w http.ResponseWriter, _ *http.Request) {
	s := fmt.Sprintf("Uptime: %s", getUptime())
	_, err := w.Write([]byte(s))
	if err != nil {
		log.Printf("baseHandlerFunc - Error: %s", err)
	}
}

func main() {
	applicationStartTime = time.Now()
	r := mux.NewRouter()
	r.HandleFunc("/", baseHandlerFunc)
	r.HandleFunc("/keys", getAllKeyHandlerFunc).Methods("GET")
	r.HandleFunc("/keys", addKeyHandlerFunc).Methods("PUT", "POST")
	r.HandleFunc("/keys/{key}", getKeyHandlerFunc).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}
