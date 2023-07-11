// auth/L402.go

package auth

import (
	"fmt"
	"log"

	models "github.com/kodylow/actually_openai/pkg/models"
	"github.com/kodylow/actually_openai/pkg/service"
)

func CheckAuthorizationHeader(reqInfo models.RequestInfo) error {
	// verify auth header exists
	if reqInfo.AuthHeader == "" {
		log.Println("No Authorization header")
		return fmt.Errorf("No Authorization header")
	}

	// Get the requestHash for checking the rune
	err := reqInfo.L402IsValid()
	if err != nil {
		log.Println("Error validating L402:", err)
		return fmt.Errorf("Invalid Authorization header")
	}

	return nil
}

func GetL402(reqInfo models.RequestInfo) (string, error) {
	// If not authorized, get msats cost for hitting this specific endpoint
	msats, err := service.MatchRequestMethodPath(reqInfo)
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
	requestHash := reqInfo.GetReqHash()
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