package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/kodylow/actually_openai/auth"
	"github.com/kodylow/actually_openai/service"
)

var APIKey string = os.Getenv("OPENAI_API_KEY")
var OpenAIEndpoint string = "https://api.openai.com/"

// Init initializes data for the handler
func Init(lnAddress string) error {
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

	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader != "L402 macaroon=macaroon:preimage" {
		// If not authorized, get invoice from lightning node
		msats := uint64(10000)
		invoice, err := service.GetInvoice(msats)
		if err != nil {
			log.Println("Error getting invoice:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// get the payment_hash out of the invoice
		paymentHash, err := service.GetPaymentHash(invoice)
		if err != nil {
			log.Println("Error getting payment hash:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Use the payment_hash in the invoice
		// to generate a Rune and restrict it to that payment hash
		master, err := auth.GetMasterRune()
		if err != nil {
			log.Println("Error getting master rune:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		restrictions := fmt.Sprintf("paymentHash=%s", paymentHash)
		restrictedRune, err := auth.GetRestrictedRune(master, restrictions)
		if err != nil {
			log.Println("Error getting restricted rune:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Successfully created restricted rune:", restrictedRune)
		token := restrictedRune

		// Send 402 Payment Required with the invoice
		w.Header().Set("WWW-Authenticate", fmt.Sprintf("L402 token=%s invoice=%s", token, invoice))
		http.Error(w, "Payment Required", http.StatusPaymentRequired)
		return
	}

	// Read the body
	bodyBytes, err := io.ReadAll(r.Body)
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
