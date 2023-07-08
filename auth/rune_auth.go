package auth

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/bolt-observer/go-runes/runes"
)

var secret []byte

func InitSecret() error {
	log.Println("Initializing Secret")

	// Read secret from environment variable
	envSecret := os.Getenv("RUNE_SECRET")

	// Convert hex encoded string secret to byte array
	var err error
	secret, err = hex.DecodeString(envSecret)
	log.Printf("Secret: %v", secret)
	if err != nil {
		log.Printf("An error occurred while initializing secret: %v", err)
		return err
	}

	log.Println("Successfully Initialized Secret")
	return nil
}

func GetMasterRune() (*runes.MasterRune, error) {
	log.Println("Creating Master Rune")

	master := runes.MustMakeMasterRune(secret)

	log.Println("Successfully Created Master Rune")
	return &master, nil
}

func GetRestrictedRune(master *runes.MasterRune, restrictions string) (*runes.Rune, error) {
	log.Println("Creating Restricted Rune")

	restricted := master.Rune.MustGetRestrictedFromString(restrictions)

	log.Println("Successfully Created Restricted Rune")
	return &restricted, nil
}