// auth/L402.go

package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kodylow/actually_openai/pkg/service"
)

func CheckAuthorizationHeader(r *http.Request) error {
	// Check Authorization header
	authHeader := r.Header.Get("Authorization")
	// verify auth header exists
	if authHeader == "" {
		log.Println("No Authorization header")
		return fmt.Errorf("No Authorization header")
	}

	// Get the requestHash for checking the rune
	reqHash := service.GetRequestHash(r)
	log.Println("Request hash:", reqHash)
	err := L402IsValid(authHeader, reqHash)
	if err != nil {
		log.Println("Error validating L402:", err)
		return fmt.Errorf("Invalid Authorization header")
	}

	return nil
}

func GetL402(r *http.Request) (string, error) {
	// If not authorized, get msats cost for hitting this specific endpoint
	msats, err := service.MatchRequestMethodPath(r)
	if err != nil {
		log.Println("Error matching request method and path for pricing:", err)
		return "", err
	}
	invoice, err := service.GetInvoice(msats)
	if err != nil {
		log.Println("Error getting invoice:", err)
		return "", err
	}
	// get the payment_hash out of the invoice
	paymentHash, err := service.GetPaymentHash(invoice)
	if err != nil {
		log.Println("Error getting payment hash:", err)
		return "", err
	}
	log.Println("Payment hash from invoice:", paymentHash)
	// get the body off the request and take the hash of it
	requestHash := service.GetRequestHash(r)
	log.Println("Calculated Request hash:", requestHash)
	// Use the payment_hash and the body_hash in the invoice
	// to generate a Rune and restrict it
	token, err := GetRestrictedRuneB64(paymentHash, requestHash)
	if err != nil {
		log.Println("Error getting restricted rune:", err)
		return "", err
	}

	l402 := fmt.Sprintf("L402 token=%s, invoice=%s", token, invoice)
	return l402, nil
}

func destructureL402AuthHeader(authHeader string) (string, string, error) {
	// Split the authHeader string by " "
	parts := strings.Split(authHeader, " ")
	// Check the parts length, it should be 2 ("L402" and "token:invoice")
	if len(parts) != 2 {
		log.Println("Invalid authorization header format destructuring L402")
		return "", "", errors.New("invalid authorization header format")
	}

	// Split the second part by ":" to get token and invoice
	tokenPreimage := strings.Split(parts[1], ":")
	log.Println("tokenPreimage:", tokenPreimage)
	// Check the tokenPreimage length, it should be 2
	if len(tokenPreimage) != 2 {
		log.Println("Invalid token:preimage format destructuring L402")
		return "", "", errors.New("invalid token:preimage format")
	}

	return tokenPreimage[0], tokenPreimage[1], nil
}

func L402IsValid(l402 string, reqHash string) error {
	// destructure off the token and preimage
	token, preimage, err := destructureL402AuthHeader(l402)
	if err != nil {
		log.Println("Error destructuring L402:", err)
		return err
	}

	// Check the token and preimage against the restrictions
	res := checkTokenRestrictions(token, preimage, reqHash)
	if !res {
		log.Println("Token doesn't match restrictions")
		return errors.New("invalid token, doesn't match restrictions")
	}

	return nil
}
