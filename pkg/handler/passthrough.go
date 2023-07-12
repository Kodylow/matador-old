package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kodylow/matador/pkg/auth"
	models "github.com/kodylow/matador/pkg/models"
	"github.com/kodylow/matador/pkg/service"
)

var APIKey string
var APIRoot string

// Init initializes data for the handler
func Init(key string, root string, lnAddress string) error {
	log.Println("LN_ADDRESS:", lnAddress)
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

	// Copy headers from the original request
	for k, v := range r.Header {
		req.Header.Set(k, v[0])
	}

	// overwrite the Authorization with the API Key
	req.Header.Set("Authorization", "Bearer "+APIKey)
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

	// Write the response back to the client
	w.Write(responseBody)

	log.Println("passthroughHandler completed")
}
