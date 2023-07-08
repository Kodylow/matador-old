package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	OpenAIEndpoint = "https://api.openai.com/"
)

var (
	APIKey string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	APIKey = os.Getenv("OPENAI_API_KEY")
}

func passthroughHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("passthroughHandler started")

	// Read the body
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new request to forward
	req, err := http.NewRequest(r.Method, OpenAIEndpoint+r.URL.Path, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Println("Error creating new forward request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy headers
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Add OpenAI API Key
	req.Header.Add("Authorization", "Bearer "+APIKey)

	// Forward the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error forwarding the request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read the response
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response back to the client
	w.Write(responseBody)

	log.Println("passthroughHandler completed")
}

func main() {
	http.HandleFunc("/", passthroughHandler)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
