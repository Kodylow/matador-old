package service

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type FunctionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Parameters  any    `json:"parameters"`
}

type ChatCompletionRequest struct {
	Model            string                         `json:"model"`
	Messages         []openai.ChatCompletionMessage `json:"messages"`
	MaxTokens        int                            `json:"max_tokens,omitempty"`
	Temperature      float32                        `json:"temperature,omitempty"`
	TopP             float32                        `json:"top_p,omitempty"`
	N                int                            `json:"n,omitempty"`
	Stream           bool                           `json:"stream,omitempty"`
	Stop             []string                       `json:"stop,omitempty"`
	PresencePenalty  float32                        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32                        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int                 `json:"logit_bias,omitempty"`
	User             string                         `json:"user,omitempty"`
	Functions        []FunctionDefinition           `json:"functions,omitempty"`
	FunctionCall     any                            `json:"function_call,omitempty"`
}

// Validate function for ChatCompletionRequest
func (self *ChatCompletionRequest) Validate() error {
	// Validate model
	if len(self.Model) == 0 {
		return fmt.Errorf("Model is required")
	}

	// Validate messages
	if len(self.Messages) == 0 {
		return fmt.Errorf("At least one message is required")
	}

	for _, msg := range self.Messages {
		if len(msg.Role) == 0 {
			return fmt.Errorf("Role is required for each message")
		}
		if msg.Role != "system" && msg.Role != "user" && msg.Role != "assistant" && msg.Role != "function" {
			return fmt.Errorf("Invalid role for message: Must be one of system, user, assistant, or function")
		}
		if len(msg.Content) == 0 && msg.FunctionCall == nil {
			return fmt.Errorf("Either content or function call is required for each message")
		}
		if msg.FunctionCall != nil && len(msg.Name) == 0 {
			return fmt.Errorf("Name is required for message with function call")
		}
	}

	// Validate n
	if self.N < 0 {
		return fmt.Errorf("Invalid number of chat completion choices: Must be non-negative")
	}

	// Validate temperature
	if self.Temperature < 0 || self.Temperature > 2 {
		return fmt.Errorf("Invalid temperature: Must be between 0 and 2")
	}

	// Validate top_p
	if self.TopP < 0 || self.TopP > 1 {
		return fmt.Errorf("Invalid top_p: Must be between 0 and 1")
	}

	// Validate presence_penalty
	if self.PresencePenalty < -2 || self.PresencePenalty > 2 {
		return fmt.Errorf("Invalid presence_penalty: Must be between -2 and 2")
	}

	// Validate frequency_penalty
	if self.FrequencyPenalty < -2 || self.FrequencyPenalty > 2 {
		return fmt.Errorf("Invalid frequency_penalty: Must be between -2 and 2")
	}

	// Validate max_tokens
	if self.MaxTokens < 0 {
		return fmt.Errorf("Invalid max_tokens: Must be non-negative")
	}

	return nil
}

// ImageGenerationBody struct
type ImageGenerationRequest struct {
	Prompt         string `json:"prompt,omitempty"`
	N              uint64 `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

// Validate function for ImageGenerationBody
func (self *ImageGenerationRequest) Validate() error {
	// Validate prompt
	if len(self.Prompt) == 0 {
		return fmt.Errorf("Prompt is required")
	}

	// Validate n
	if self.N <= 0 {
		self.N = 1 // Default value is 1
	} else if self.N > 10 {
		return fmt.Errorf("Invalid number of images to generate: Must be between 1 and 10")
	}

	// Validate size
	if self.Size == "" {
		self.Size = "1024x1024" // Default value is 1024x1024
	}

	// Validate response format
	if self.ResponseFormat == "" {
		self.ResponseFormat = "url" // Default value is url
	} else if self.ResponseFormat != "url" && self.ResponseFormat != "b64_json" {
		return fmt.Errorf("Invalid Image Generation Response Format: Must be one of url or b64_json")
	}

	return nil
}
