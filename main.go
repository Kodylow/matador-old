package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/niftynei/glightning/glightning"
)

const (
	OpenAIEndpoint = "https://api.openai.com/"
)

var (
	APIKey    string
	Lightning *glightning.Lightning
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	APIKey = os.Getenv("OPENAI_API_KEY")

	// Start up the lightning node
	Lightning = glightning.NewLightning()
	Lightning.StartUp("lightning-rpc", "/tmp/clight-1")
}

// passthroughHandler forwards the request to the OpenAI API
func passthroughHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("passthroughHandler started")

	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "L402 macaroon=macaroon:preimage" {
		// If not authorized, get invoice from lightning node
		satoshi := uint64(10000)
		invoiceLabel := "ayc"
		invoice, err := Lightning.CreateInvoice(satoshi, invoiceLabel, "desc", uint32(5), nil, "", false)
		if err != nil {
			log.Println("Error creating invoice:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send 402 Payment Required with the invoice
		w.Header().Set("WWW-Authenticate", fmt.Sprintf("L402 macaroon=macaroon invoice=%s", invoice.Bolt11))
		http.Error(w, "Payment Required", http.StatusPaymentRequired)
		return
	}

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