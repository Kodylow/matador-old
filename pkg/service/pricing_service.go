package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/kodylow/actually_openai/pkg/utils"
)

// Endpoint structure
type Endpoint struct {
  MethodPath string
  PriceFunc func(r *http.Request) (uint64, error)
}

func getMsats(r *http.Request) (uint64, error) {
  return 10000, nil
}

type ImageGenerationBody struct {
	Prompt         string `json:"prompt"`
	N              uint64    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
	User           string `json:"user"`
}

func getMsatsImageGenerations(r *http.Request) (uint64, error) {
	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		return 0, fmt.Errorf("Error reading request body: %w", err)
	}
	// Reset the request body to the original state
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Get the body of the request as a ImageGenerationBody
	var body ImageGenerationBody
	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&body)
	if err != nil {
		log.Println("Error decoding ImageGeneration request body:", err)
		return 0, fmt.Errorf("Error decoding ImageGeneration request body: %w", err)
	}

	// Validate prompt
	if len(body.Prompt) == 0 {
		log.Println("Prompt is required")
		return 0, fmt.Errorf("Prompt is required")
	}

	// Validate n
	if body.N <= 0 {
		body.N = 1 // Default value is 1
	} else if body.N > 10 {
		log.Println("Invalid number of images to generate:", body.N)
		return 0, fmt.Errorf("Invalid number of images to generate: Must be between 1 and 10")
	}

	// Validate size
	if body.Size == "" {
		body.Size = "1024x1024" // Default value is 1024x1024
	}

	var cents uint64
	switch body.Size {
	case "256x256":
		cents = body.N * 16 // $0.016 per image
	case "512x512":
		cents = body.N * 18 // $0.018 per image
	default:
		cents = body.N * 20 // $0.020 per image
	}

	log.Println("Cents:", cents)
	// convert cents to msats
	msats := utils.CentsToMsats(cents)

	// Validate response format
	if body.ResponseFormat == "" {
		body.ResponseFormat = "url" // Default value is url
	} else if body.ResponseFormat != "url" && body.ResponseFormat != "b64_json" {
		log.Println("Invalid Image Generation Response Format:", body.ResponseFormat)
		return 0, fmt.Errorf("Invalid Image Generation Response Format: Must be one of url or b64_json")
	}

	return msats, nil
}


// List of endpoints
var endpoints = []Endpoint{
  {"GET /v1/models", getMsats},
  {"GET /v1/models/", getMsats},
  {"POST /v1/chat/completions", getMsats},
  {"POST /v1/completions", getMsats},
  {"POST /v1/images/generations", getMsatsImageGenerations},
  {"POST /v1/images/edits", getMsats},
  {"POST /v1/images/variations", getMsats},
  {"POST /v1/embeddings", getMsats},
  {"POST /v1/audio/transcriptions", getMsats},
  {"POST /v1/audio/translations", getMsats},
  {"GET /v1/files", getMsats},
  {"POST /v1/files", getMsats},
  {"DELETE /v1/files/", getMsats},
  {"GET /v1/files/", getMsats},
  {"GET /v1/files/content", getMsats},
  {"POST /v1/fine-tunes", getMsats},
  {"GET /v1/fine-tunes", getMsats},
  {"GET /v1/fine-tunes/", getMsats},
  {"POST /v1/fine-tunes/cancel", getMsats},
  {"GET /v1/fine-tunes/events", getMsats},
  {"DELETE /v1/models/", getMsats},
  {"POST /v1/moderations", getMsats},
}

// MatchRequestMethodPath function
func MatchRequestMethodPath(r *http.Request) (uint64, error) {
  // Make MethodPath from request
  methodPath := strings.ToUpper(r.Method) + " " + r.URL.Path
  for _, endpoint := range endpoints {
	// Match MethodPath with endpoint.MethodPath
	matched, err := regexp.MatchString(endpoint.MethodPath, methodPath)
	if err != nil {
	  log.Println("Error matching string:", err)
	}
    if matched {
	  log.Println("Matched:", endpoint.MethodPath)
      return endpoint.PriceFunc(r)
    }
  }
  return 21000, nil
}