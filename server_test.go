package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetKeyNotFound(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/keys/{key}", GetKeyHandlerFunc)
	req, err := http.NewRequest("GET", "/keys/key1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v expected %v",
			status, http.StatusNotFound)
	}

	if rr.Body.String() != "" {
		t.Errorf("handler returned unexpected body: got %v expected %v",
			rr.Body.String(), "")
	}
}

func TestHandlerEmptyGetAll(t *testing.T) {
	req, err := http.NewRequest("GET", "/keys", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllKeyHandlerFunc)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}

	expected := `[]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v expected %v",
			rr.Body.String(), expected)
	}
}

func TestHandlerGetAll(t *testing.T) {
	numPairs := 10
	expected := make([]string, numPairs)

	for i := 0; i < numPairs; i++ {
		tKey := fmt.Sprintf("Key%d", i)
		tVal := fmt.Sprintf("Value%d", i)
		expected[i] = fmt.Sprintf(`{"key":"Key%d","value":"Value%d"}`, i, i)

		err := Put(tKey, tVal)
		if err != nil {
			t.Errorf("Got error when attempting to Put() Key: %s and Value %s", tKey, tVal)
		}
	}

	req, err := http.NewRequest("GET", "/keys", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllKeyHandlerFunc)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, expected %v",
			status, http.StatusOK)
	}

	resBdy := rr.Body.String()
	for _, res := range expected {
		if !strings.Contains(resBdy, res) {
			t.Errorf("Res Body does not contain expected value: %s", res)
		}
	}
}

func TestHandlerFoundKey(t *testing.T) {
	key := "TestKey"
	value := "TestVal"

	err := Put(key, value)
	if err != nil {
		t.Errorf("Failure adding test pair")
	}

	path := fmt.Sprintf("/keys/%s", key)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/keys/{key}", GetKeyHandlerFunc)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error - Expected StatusCode %d, got %d", http.StatusOK, rr.Code)
	}

	resBdy := rr.Body.String()
	if resBdy != value+"\n" {
		t.Errorf("Incorrect value returned for key: %s, got %s, expected %s", key, resBdy, value)
	}
}

func TestPostHandlerFoundKey(t *testing.T) {
	key := "TestKey2"
	value := "TestVal2"

	router := mux.NewRouter()
	router.HandleFunc("/keys/{key}", GetKeyHandlerFunc)
	router.HandleFunc("/keys", AddKeyHandlerFunc).Methods("POST")

	reqBody := fmt.Sprintf(`{"key": "%s", "value": "%s"}`, key, value)
	reqCreate, err := http.NewRequest("POST", "/keys", bytes.NewBuffer([]byte(reqBody)))
	reqCreate.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("Got Error when POST'ing: %v", err.Error())
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, reqCreate)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Error - Expected StatusCode Created %d, got %d", http.StatusCreated, rr.Code)
	}

	reqGet, err := http.NewRequest("GET", fmt.Sprintf("/keys/%s", key), nil)
	if err != nil {
		t.Fatalf("Error: %v during request to get %s", err.Error(), key)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, reqGet)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error - Expected StatusCode OK %d, got %d", http.StatusOK, rr.Code)
	}

	resBdy := rr.Body.String()
	if resBdy != value+"\n" {
		t.Errorf("Incorrect value returned for key: %s, got %s, expected %s", key, resBdy, value)
	}
}

func TestPostHandlerUpdate(t *testing.T) {
	key := "TestKey3"
	value := "TestVal3"
	updateVal := "Update3"

	router := mux.NewRouter()
	router.HandleFunc("/keys/{key}", GetKeyHandlerFunc)
	router.HandleFunc("/keys", AddKeyHandlerFunc).Methods("POST", "PUT")

	reqBody := fmt.Sprintf(`{"key":"%s","value":"%s"}`, key, value)
	reqCreate, err := http.NewRequest("POST", "/keys", bytes.NewBuffer([]byte(reqBody)))
	reqCreate.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("Got Error when POST'ing: %v", err.Error())
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, reqCreate)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Error - Expected StatusCode Created %d, got %d", http.StatusCreated, rr.Code)
	}

	// test PUT update
	reqBody = fmt.Sprintf(`{"key":"%s","value":"%s"}`, key, updateVal)
	reqPut, err := http.NewRequest("PUT", "/keys", bytes.NewBuffer([]byte(reqBody)))
	reqPut.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatalf("Error: %v during request to PUT %s", err.Error(), key)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, reqPut)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error - Expected StatusCode OK %d, got %d", http.StatusOK, rr.Code)
	}

	resBdy := rr.Body.String()
	if resBdy != reqBody {
		t.Errorf("Incorrect value returned for key: %s, got %s, expected %s", key, resBdy, reqBody)
	}

	// test GET - the updated value is returned
	reqGet, err := http.NewRequest("GET", fmt.Sprintf("/keys/%s", key), nil)
	if err != nil {
		t.Fatalf("Error: %v during request to get %s", err.Error(), key)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, reqGet)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Error - Expected StatusCode OK %d, got %d", http.StatusOK, rr.Code)
	}

	resBdy = rr.Body.String()
	if resBdy != updateVal+"\n" {
		t.Errorf("Incorrect value returned for key: %s, got %s, expected %s", key, resBdy, updateVal)
	}
}
