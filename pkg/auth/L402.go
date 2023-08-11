// auth/L402.go

package auth

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/kodylow/matador/pkg/database"
	models "github.com/kodylow/matador/pkg/models"
	"github.com/kodylow/matador/pkg/service"
)

// ExtractToken extracts the token from the Authorization header
func ExtractToken(authHeader string) (string, error) {
	log.Println("Extracting token from Authorization header:", authHeader)
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "L402" {
		return "", fmt.Errorf("Invalid Authorization header format: should be L402 token:preimage")
	}

	params := strings.Split(parts[1], ":")
	if len(params) != 2 {
		return "", fmt.Errorf("Invalid Authorization header format, should be L402 token:preimage")
	}

	return params[0], nil
}

func CheckAuthorizationHeader(reqInfo models.RequestInfo) error {
	// verify auth header exists
	if reqInfo.AuthHeader == "" {
		log.Println("No Authorization header")
		return fmt.Errorf("No Authorization header")
	}

	// Extract the token from the Authorization header
	token, err := ExtractToken(reqInfo.AuthHeader)
	if err != nil {
		log.Println("Error extracting token from Authorization header:", err)
		return fmt.Errorf("Invalid Authorization header")
	}

	// Check if the token exists in the database
	spent, err := database.GetToken(token)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Token not found in the database")
			return fmt.Errorf("Invalid token")
		}
		log.Println("Error retrieving token from the database:", err)
		return fmt.Errorf("Invalid Authorization header")
	}

	// Check if the token has been spent
	if spent {
		log.Println("Token has already been spent")
		return fmt.Errorf("Token has already been spent")
	}

	// Mark the token as spent
	err = database.UpdateToken(token, true)
	if err != nil {
		log.Println("Error marking token as spent:", err)
		return fmt.Errorf("Invalid Authorization header")
	}

	// Get the requestHash for checking the rune
	err = reqInfo.L402IsValid()
	if err != nil {
		log.Println("Error validating L402:", err)
		return fmt.Errorf("Invalid Authorization header")
	}

	return nil
}

// GetL402 returns the L402 authentication header for the given request
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
	// get the body off the request and take the hash of it
	requestHash := reqInfo.GetReqHash()
	// Use the payment_hash and the body_hash in the invoice
	// to generate a Rune and restrict it
	token, err := GetRestrictedRuneB64(paymentHash, requestHash)
	if err != nil {
		log.Println("Error getting restricted rune:", err)
		return "", err
	}

	// Store the issued token in the database as unspent
	err = database.AddToken(token)
	if err != nil {
		log.Println("Error storing issued token in the database:", err)
		return "", err
	}

	l402 := fmt.Sprintf("L402 token=%s, invoice=%s", token, invoice)
	return l402, nil
}
