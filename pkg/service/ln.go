package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	ln "github.com/nbd-wtf/ln-decodepay"
)

type LnAddressResponse struct {
	Callback    string `json:"callback"`
	MinSendable uint64 `json:"minSendable"`
	MaxSendable uint64 `json:"maxSendable"`
}

type LnCallbackResponse struct {
	Pr     string `json:"pr"`
	Status string `json:"status"`
	Verify string `json:"verify"`
}

var LnAddr LnAddressResponse

// GetCallback gets the callback from the lightning address
func GetCallback(lnAddress string) (LnAddressResponse, error) {
	parts := strings.Split(lnAddress, "@")
	if len(parts) != 2 {
		return LnAddressResponse{}, fmt.Errorf("invalid lnAddress: %s", lnAddress)
	}
	username := parts[0]
	domain := parts[1]
	url := fmt.Sprintf("https://%s/.well-known/lnurlp/%s", domain, username)
	resp, err := http.Get(url)
	if err != nil {
		return LnAddressResponse{}, err
	}
	if resp.StatusCode != 200 {
		return LnAddressResponse{}, fmt.Errorf("invalid lnAddress: %s", lnAddress)
	}
	var lnAddrResp LnAddressResponse
	err = json.NewDecoder(resp.Body).Decode(&lnAddrResp)
	if err != nil {
		log.Fatal("Error decoding callback response:", err)
		return LnAddressResponse{}, err
	}

	return lnAddrResp, nil
}

// GetInvoice gets an invoice from the lightning callback
func GetInvoice(msats uint64) (string, error) {
	if msats > LnAddr.MaxSendable || msats < LnAddr.MinSendable {
		return "", fmt.Errorf("%d msats not in sendable range of %d - %d:", msats, LnAddr.MinSendable, LnAddr.MaxSendable)
	}
	resp, err := http.Get(fmt.Sprintf("%s?amount=%d", LnAddr.Callback, msats))
	if err != nil {
		log.Println("Error getting invoice:", err)
		return "", err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	var lnCallbackResp LnCallbackResponse
	err = json.Unmarshal(bodyBytes, &lnCallbackResp)
	if err != nil {
		log.Println("Error decoding callback response:", err)
		return "", err
	}

	log.Println("Verify Url:", lnCallbackResp.Verify)

	return lnCallbackResp.Pr, nil
}

func GetPaymentHash(invoice string) (string, error) {
	bolt11, err := ln.Decodepay(invoice)
	if err != nil {
		return "", fmt.Errorf("error decoding invoice: %w", err)
	}

	return bolt11.PaymentHash, nil
}
