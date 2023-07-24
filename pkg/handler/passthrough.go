package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kodylow/matador/pkg/auth"
	models "github.com/kodylow/matador/pkg/models"
	"github.com/kodylow/matador/pkg/service"
)

var APIKey string
var APIRoot string

// Init initializes data for the handler
func Init(key string, root string, lnAddress string) error {
	APIKey = key
	APIRoot = root
	var err error
	service.LnAddr, err = service.GetCallback(lnAddress)
	if err != nil {
		return fmt.Errorf("error getting lnaddress callback: %w", err)
	}
	return nil
}

// PassthroughHandler forwards the request to the OpenAI API
func PassthroughHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("passthroughHandler started")

	// (Your code remains the same up to here)

	// Forward the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error forwarding the request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers from the response
	for name, values := range resp.Header {
		w.Header()[name] = values
	}

	// Set the status code on the response writer to the status code of the response
	w.WriteHeader(resp.StatusCode)

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Response from OpenAI:", string(responseBody))

	// Write the response back to the client
	w.Write(responseBody)

	log.Println("passthroughHandler completed")
}

// ModelsHandler returns a list of all models
func ModelsHandler(w http.ResponseWriter, r *http.Request) {
	// Hardcode the path and method
	r.URL.Path = "/v1/models"
	r.Method = "GET"
	PassthroughHandler(w, r)
}

// SpecificModelHandler returns a specific model
func SpecificModelHandler(w http.ResponseWriter, r *http.Request) {
	// Get the model ID from the URL
	id := mux.Vars(r)["id"]
	// Hardcode the path and method
	r.URL.Path = "/v1/models/" + id
	r.Method = "GET"
	PassthroughHandler(w, r)
}

// ChatCompletionsHandler handles chat completions
func ChatCompletionsHandler(w http.ResponseWriter, r *http.Request) {
	// Hardcode the path and method
	r.URL.Path = "/v1/chat/completions"
	r.Method = "POST"
	PassthroughHandler(w, r)
}

// ImagesGenerationsHandler handles image generations
func ImagesGenerationsHandler(w http.ResponseWriter, r *http.Request) {
	// Hardcode the path and method
	r.URL.Path = "/v1/images/generations"
	r.Method = "POST"
	PassthroughHandler(w, r)
}

// EmbeddingsHandler handles embeddings
func EmbeddingsHandler(w http.ResponseWriter, r *http.Request) {
	// Hardcode the path and method
	r.URL.Path = "/v1/embeddings"
	r.Method = "POST"
	PassthroughHandler(w, r)
}
