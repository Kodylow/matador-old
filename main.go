package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kodylow/actually_openai/auth" // import the auth package
	"github.com/kodylow/actually_openai/handler"
)

var APIKey string

const OpenAIEndpoint = "https://api.openai.com/"

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	APIKey = os.Getenv("OPENAI_API_KEY")
	err = handler.Init(os.Getenv("LN_ADDRESS"))
	if err != nil {
		log.Fatal("Error initializing handlers: ", err)
	}

	// Initialize the secret
	err = auth.InitSecret() // read the secret from the RUNE_SECRET environment variable
	if err != nil {
		log.Fatal("Error initializing secret: ", err)
	}
}

func main() {
	http.HandleFunc("/", handler.PassthroughHandler)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}