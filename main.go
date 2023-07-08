package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	OpenAIEndpoint = "https://api.openai.com/v1/"
)

var (
	APIKey string
)

func init() {
	APIKey = os.Getenv("OPENAI_API_KEY")  // Read API key from environment variable
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read the body
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a new request to forward
		req, err := http.NewRequest(r.Method, OpenAIEndpoint + r.URL.Path, bytes.NewBuffer(bodyBytes))
		if err != nil {
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Read the response
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the response back to the client
		w.Write(responseBody)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
