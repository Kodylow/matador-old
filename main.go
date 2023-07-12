package main

import (
	"log"
	"net/http"
	"os"

	// "github.com/joho/godotenv"
	"github.com/kodylow/renegade/pkg/auth"
	"github.com/kodylow/renegade/pkg/database"
	"github.com/kodylow/renegade/pkg/handler"
)

func init() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	err := handler.Init(os.Getenv("API_KEY"), os.Getenv("API_ROOT"), os.Getenv("LN_ADDRESS"))
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
	http.HandleFunc("/", handler.RootHandler)
	http.HandleFunc("/v1/", handler.PassthroughHandler)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
