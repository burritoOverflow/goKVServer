package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var applicationStartTime time.Time

type KeyValEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

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
	w.Write([]byte(keyRes + "\n"))
}

func addKeyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// add new key
	case http.MethodPost:
		// if exists add anyway
		var kvEntry KeyValEntry
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
			w.Write([]byte("Error"))
		}
		w.Write(jsonKv)
	}
}

func baseHandlerFunc(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("Uptime: %s", getUptime())
	w.Write([]byte(s))
}

func main() {
	applicationStartTime = time.Now()
	r := mux.NewRouter()
	r.HandleFunc("/", baseHandlerFunc)
	r.HandleFunc("/keys", addKeyHandlerFunc).Methods("PUT", "POST")
	r.HandleFunc("/keys/{key}", getKeyHandlerFunc).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
