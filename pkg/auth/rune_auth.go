package auth

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/bolt-observer/go-runes/runes"
	"github.com/kodylow/actually_openai/service"
)

var secret []byte

// InitSecret initializes the secret from the environment variable
func InitSecret() error {
	log.Println("Initializing Secret")

	// Read secret from environment variable
	envSecret := os.Getenv("RUNE_SECRET")

	// Convert hex encoded string secret to byte array
	var err error
	secret, err = hex.DecodeString(envSecret)
	if err != nil {
		log.Printf("An error occurred while initializing secret: %v", err)
		return err
	}

	log.Println("Successfully Initialized Secret")
	return nil
}

// GetRestrictedRune returns a rune reqstricted to paymentHash and requestHash from a master rune
func GetRestrictedRuneB64(paymentHash string, requestHash string) (string, error) {
	master := runes.MustMakeMasterRune(secret)
	restrictedRune := master.Rune.MustGetRestrictedFromString("paymentHash=" + paymentHash + "&" + "requestHash=" + requestHash)
	restrictedRuneB64 := restrictedRune.ToBase64()
	return restrictedRuneB64, nil
}

// L402IsValid checks if the given rune is valid
func checkTokenRestrictions(runeB64 string, preimage string, reqHash string) bool {
    // hash the preimage to get the paymentHash
    hash := service.Sha256Hash(preimage)
	log.Println("Payment Hash Calculated from Preimage:", hash)
	// get the master rune from the secret
	master := runes.MustMakeMasterRune(secret)
	log.Println("Master Rune:", master.Rune.ToBase64())
    // decode the given rune from base64
    restrictedRune := runes.MustGetFromBase64(runeB64)
	log.Println("Restricted Rune:", restrictedRune)
    // create map with the values to evaluate
    values := map[string]any{
        "paymentHash": hash,
        "requestHash": reqHash,
    }

    // evaluate the rune to check if the given hashes match the restrictions
    err := master.Check(&restrictedRune, values)
	if err != nil {
		log.Println("Error checking rune:", err)
		return false
	}
    return true
}
