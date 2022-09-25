package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
