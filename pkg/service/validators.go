package service

import (
	"fmt"
)

// ImageGenerationBody struct
type ImageGenerationBody struct {
    Prompt         string `json:"prompt"`
    N              uint64 `json:"n"`
    Size           string `json:"size"`
    ResponseFormat string `json:"response_format"`
    User           string `json:"user"`
}

// Validate function for ImageGenerationBody
func (body *ImageGenerationBody) Validate() error {
    // Validate prompt
    if len(body.Prompt) == 0 {
        return fmt.Errorf("Prompt is required")
    }

    // Validate n
    if body.N <= 0 {
        body.N = 1 // Default value is 1
    } else if body.N > 10 {
        return fmt.Errorf("Invalid number of images to generate: Must be between 1 and 10")
    }

    // Validate size
    if body.Size == "" {
        body.Size = "1024x1024" // Default value is 1024x1024
    }

    // Validate response format
    if body.ResponseFormat == "" {
        body.ResponseFormat = "url" // Default value is url
    } else if body.ResponseFormat != "url" && body.ResponseFormat != "b64_json" {
        return fmt.Errorf("Invalid Image Generation Response Format: Must be one of url or b64_json")
    }

    return nil
}
