package main

import (
	"github.com/kodylow/actually_openai/handler"
	"log"
	"net/http"
	"os"
)

var APIKey string

const OpenAIEndpoint = "https://api.openai.com/"

func init() {
	APIKey = os.Getenv("OPENAI_API_KEY")
	err := handler.Init(os.Getenv("LN_ADDRESS"))
	if err != nil {
		log.Fatal("Error initializing handlers: ", err)
	}
}

func main() {
	http.HandleFunc("/", handler.PassthroughHandler)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
