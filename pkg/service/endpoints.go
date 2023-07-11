package service

import (
	"path"

	models "github.com/kodylow/actually_openai/pkg/models"
)

// Endpoint structure
type Endpoint struct {
	Method    string
	Path      string
	PriceFunc func(models.RequestInfo) (uint64, error)
}

// List of endpoints
var supportedEndpoints = []Endpoint{
	{"GET", "/v1/models", get1000Msats},                       // done
	{"GET", "/v1/models/", get1000Msats},                      // done
	{"POST", "/v1/chat/completions", getMsatsChatCompletions}, //done
	// {"POST", "/v1/completions", getMsats},
	{"POST", "/v1/images/generations", getMsatsImageGenerations}, // done
	// {"POST", "/v1/images/edits", getMsats},
	// {"POST", "/v1/images/variations", getMsats},
	// {"POST", "/v1/embeddings", getMsats},
	// {"POST", "/v1/audio/transcriptions", getMsats},
	// {"POST", "/v1/audio/translations", getMsats},
	//   {"GET /v1/files", getMsats},
	//   {"POST /v1/files", getMsats},
	//   {"DELETE /v1/files/", getMsats},
	//   {"GET /v1/files/", getMsats},
	//   {"GET /v1/files/content", getMsats},
	//   {"POST /v1/fine-tunes", getMsats},
	//   {"GET /v1/fine-tunes", getMsats},
	//   {"GET /v1/fine-tunes/", getMsats},
	//   {"POST /v1/fine-tunes/cancel", getMsats},
	//   {"GET /v1/fine-tunes/events", getMsats},
	//   {"DELETE /v1/models/", getMsats},
	//   {"POST /v1/moderations", getMsats},
}

// MatchRequestMethodPath function
func MatchRequestMethodPath(reqInfo models.RequestInfo) (uint64, error) {
	for _, endpoint := range supportedEndpoints {
		matched_path, err := path.Match(endpoint.Path, reqInfo.Path)
		if err != nil {
			return 0, err
		}
		if reqInfo.Method == endpoint.Method && matched_path {
			return endpoint.PriceFunc(reqInfo)
		}
	}
	return 21000, nil // Default pricing
}
