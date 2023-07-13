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

	// Check if the requested endpoint is supported
	reqPath := r.URL.Path
	reqMethod := r.Method
	var isSupported bool
	for _, endpoint := range service.SupportedEndpoints {
		if reqPath == endpoint.Path && reqMethod == endpoint.Method {
			isSupported = true
			break
		}
	}

	if !isSupported {
		http.Error(w, "Endpoint not supported", http.StatusNotFound)
		return
	}

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Reset the body to its original state

	// Create a RequestInfo
	reqInfo := models.RequestInfo{
		AuthHeader: r.Header.Get("Authorization"),
		Method:     r.Method,
		Path:       r.URL.Path,
		Body:       body,
	}

	err = auth.CheckAuthorizationHeader(reqInfo)
	if err != nil {
		log.Println("Unauthorized, payment required")
		l402, err := auth.GetL402(reqInfo)
		if err != nil {
			log.Println("Error getting L402:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Send 402 Payment Required with the invoice
		w.Header().Set("WWW-Authenticate", l402)
		http.Error(w, "Payment Required", http.StatusPaymentRequired)
		return
	}

	// Create a new request to forward
	req, err := http.NewRequest(reqInfo.Method, APIRoot+reqInfo.Path, r.Body)
	if err != nil {
		log.Println("Error creating new forward request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// overwrite the Authorization with the API Key
	req.Header.Set("Authorization", "Bearer "+APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	log.Println("Forwarding request to OpenAI API")
	log.Println("Request:", req)
	// Forward the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error forwarding the request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

// AudioTranscriptionsHandler handles audio transcriptions
func AudioTranscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	// Hardcode the path and method
	r.URL.Path = "/v1/audio/transcriptions"
	r.Method = "POST"
	PassthroughHandler(w, r)
}
