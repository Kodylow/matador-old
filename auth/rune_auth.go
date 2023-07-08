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

// GetMasterRune returns a master rune from the secret
func GetMasterRune() (*runes.MasterRune, error) {
	master := runes.MustMakeMasterRune(secret)
	return &master, nil
}

// GetRestrictedRune returns a restricted rune from a master rune
func GetRestrictedRune(master *runes.MasterRune, restrictions string) (*runes.Rune, error) {
	restricted := master.Rune.MustGetRestrictedFromString(restrictions)
	return &restricted, nil
}