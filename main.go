package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	// "github.com/joho/godotenv"
	// "github.com/niftynei/glightning/glightning"
)

const (
	OpenAIEndpoint = "https://api.openai.com/"
)

var (
	APIKey string
	// Lightning *glightning.Lightning
	LnAddr LnAddressResponse
)

type LnAddressResponse struct {
	Callback    string `json:"callback"`
	MinSendable int64  `json:"minSendable"`
	MaxSendable int64  `json:"maxSendable"`
}

type LnCallbackResponse struct {
	Pr     string `json:"pr"`
	Status string `json:"status"`
}

func getCallback(lnAddress string) (LnAddressResponse, error) {
	// split lnAddress into <username>@<domain>
	parts := strings.Split(lnAddress, "@")
	if len(parts) != 2 {
		return LnAddressResponse{}, fmt.Errorf("invalid lnAddress: %s", lnAddress)
	}
	username := parts[0]
	domain := parts[1]
	// construct the lnAddress check url
	url := fmt.Sprintf("https://%s/.well-known/lnurlp/%s", domain, username)
	// make the request
	resp, err := http.Get(url)
	if err != nil {
		return LnAddressResponse{}, err
	}
	if resp.StatusCode != 200 {
		return LnAddressResponse{}, fmt.Errorf("invalid lnAddress: %s", lnAddress)
	}
	// pull the callback off the response body
	var lnAddrResp LnAddressResponse
	err = json.NewDecoder(resp.Body).Decode(&lnAddrResp)
	if err != nil {
		log.Fatal("Error decoding callback response:", err)
	}

	return lnAddrResp, nil
}

func getInvoice(msats int64) (string, error) {
	if msats > LnAddr.MaxSendable || msats < LnAddr.MinSendable {
		return "", fmt.Errorf("%d msats not in sendable range of %s - %s:", msats, LnAddr.MinSendable, LnAddr.MaxSendable)
	}
	resp, err := http.Get(fmt.Sprintf("%s?amount=%s", LnAddr.Callback, msats))
	if err != nil {
		log.Println("Error getting invoice:", err)
		return "", err
	}
	// parse out the LnCallbackResponse
	var lnCallbackResp LnCallbackResponse
	err = json.NewDecoder(resp.Body).Decode(&lnCallbackResp)
	if err != nil || lnCallbackResp.Status != "OK" {
		log.Println("Error decoding callback response:", err)
		return "", err
	}
	// return the payment request
	return lnCallbackResp.Pr, nil
}

func init() {
	APIKey = os.Getenv("OPENAI_API_KEY")
	LnAddress := os.Getenv("LN_ADDRESS")

	// Get callback URL from lightning address
	LnAddr, err := getCallback(LnAddress)
	if err != nil {
		log.Fatal("Error getting lnaddress callback:", err)
	}
}

// passthroughHandler forwards the request to the OpenAI API
func passthroughHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("passthroughHandler started")

	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "L402 macaroon=macaroon:preimage" {
		// If not authorized, get invoice from lightning node
		msats := int64(10000)
		invoice, err := getInvoice(msats)
        if err != nil {
            log.Println("Error getting invoice:", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

		// Send 402 Payment Required with the invoice
		w.Header().Set("WWW-Authenticate", fmt.Sprintf("L402 macaroon=macaroon invoice=%s", invoice))
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
