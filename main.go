package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	// "github.com/joho/godotenv"
	"github.com/kodylow/matador/pkg/auth"
	"github.com/kodylow/matador/pkg/database"
	"github.com/kodylow/matador/pkg/handler"
	"github.com/rs/cors"
)

func init() {
	// read in .env
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file: ", err)
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
	router := mux.NewRouter()

	// Root handler
	// router.HandleFunc("/", handler.RootHandler

	router.HandleFunc("/v1/models", handler.ModelsHandler).Methods("GET")
	router.HandleFunc("/v1/models/", handler.SpecificModelHandler).Methods("GET")
	router.HandleFunc("/v1/chat/completions", handler.ChatCompletionsHandler).Methods("POST")
	router.HandleFunc("/v1/images/generations", handler.ImagesGenerationsHandler).Methods("POST")
	router.HandleFunc("/v1/embeddings", handler.EmbeddingsHandler).Methods("POST")

	// Add demo routes
	demoRouter := router.PathPrefix("/demo").Subrouter()
	demoRouter.HandleFunc("/", handler.DemoIndexHandler)
	demoRouter.HandleFunc("/chat", handler.DemoChatHandler)
	demoRouter.HandleFunc("/images", handler.DemoImageHandler)
	demoRouter.HandleFunc("/embeddings", handler.DemoEmbeddingsHandler)

	// add static file server
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	router.PathPrefix("/demo/static/").Handler(http.StripPrefix("/demo/static/", http.FileServer(http.Dir("./static/"))))

	router.HandleFunc("/", handler.DemoIndexHandler)

	// setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // change this to the domains you want to al
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"}, // change this to the headers you want to allow
		ExposedHeaders:   []string{"*"}, // change this to the headers you want to expose
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Println("Gorilla CORS Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
