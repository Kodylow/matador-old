package auth

import (
	"crypto/rand"
	"log"

	"github.com/bolt-observer/go-runes/runes"
)

var secret []byte

func InitSecret() error {
	log.Println("Initializing Secret")

	secret = make([]byte, 55)
	_, err := rand.Read(secret)
	if err != nil {
		log.Printf("An error occurred while initializing secret: %v", err)
		return err
	}

	log.Println("Successfully Initialized Secret")
	return nil
}

func GetMasterRune() (*runes.Rune, error) {
	log.Println("Creating Master Rune")

	master := runes.MustMakeMasterRune(secret)

	log.Println("Successfully Created Master Rune")
	return &master.Rune, nil
}

func GetRestrictedRune(master *runes.Rune, restrictions string) (*runes.Rune, error) {
	log.Println("Creating Restricted Rune")

	restricted := master.MustGetRestrictedFromString(restrictions)

	log.Println("Successfully Created Restricted Rune")
	return &restricted.Rune, nil
}
