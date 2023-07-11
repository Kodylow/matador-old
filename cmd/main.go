package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kodylow/actually_openai/pkg/auth" // import the auth package
	"github.com/kodylow/actually_openai/pkg/database"
	"github.com/kodylow/actually_openai/pkg/handler"
)

var APIKey string

const APIRoot = "https://api.openai.com/"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = handler.Init(os.Getenv("API_KEY"), os.Getenv("LN_ADDRESS"))
	if err != nil {
		log.Fatal("Error initializing environment variables for handlers: ", err)
	}

	err = database.InitDatabase()
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}

	// Initialize the secret
	err = auth.InitSecret() // read the secret from the RUNE_SECRET environment variable
	if err != nil {
		log.Fatal("Error initializing secret for server side tokens/runes: ", err)
	}
}

func main() {
	http.HandleFunc("/", handler.PassthroughHandler)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
