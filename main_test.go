package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPassthrough(t *testing.T) {
	log.Println("TestPassthrough started")

	// Create a request to pass to our handler.
	reqBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant."},
			{"role": "user", "content": "Hello!"},
		},
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("Error marshaling request body:", err)
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/chat/completions", bytes.NewBuffer(reqBytes))
	if err != nil {
		log.Println("Error creating new request:", err)
		t.Fatal(err)
	}

	// We create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()

	// Call the passthrough handler
	handler := http.HandlerFunc(passthroughHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		log.Printf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	log.Println("TestPassthrough completed")
}
