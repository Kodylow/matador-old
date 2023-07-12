package auth

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/bolt-observer/go-runes/runes"
)

var secret []byte

// InitSecret initializes the secret from the environment variable
func InitSecret() error {
	// Read secret from environment variable
	envSecret := os.Getenv("RUNE_SECRET")

	// Convert hex encoded string secret to byte array
	var err error
	secret, err = hex.DecodeString(envSecret)
	if err != nil {
		log.Printf("An error occurred while initializing secret: %v", err)
		return err
	}
	return nil
}

// GetRestrictedRune returns a rune reqstricted to paymentHash and requestHash from a master rune
func GetRestrictedRuneB64(paymentHash string, requestHash string) (string, error) {
	master := runes.MustMakeMasterRune(secret)
	restrictedRune := master.Rune.MustGetRestrictedFromString("paymentHash=" + paymentHash + "&" + "requestHash=" + requestHash)
	restrictedRuneB64 := restrictedRune.ToBase64()
	return restrictedRuneB64, nil
}
