package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

const BTCPRICE float64 = 25000 // $25,000.00 per bitcoin

func DollarsToMsats(dollars float64) uint64 {
	btc := dollars / (BTCPRICE) // Convert cents to BTC
	msats := btc * 100000000000 // Convert BTC to msats
	return uint64(msats)
}

func ModelPricingDollarsPer1KTokens(model string) (float64, error) {
	switch model {
	case "gpt-4", "gpt-4-0613":
		return 0.06, nil
	case "gpt-4-32k", "gpt-4-32k-0613":
		return 0.12, nil
	case "gpt-3.5-turbo", "gpt-3.5-turbo-0613":
		return 0.002, nil
	case "gpt-3.5-turbo-16k", "gpt-3.5-turbo-16k-0613":
		return 0.004, nil
	default:
		return 0, fmt.Errorf("Unknown model: %s", model)
	}
}

func TokensToMsats(tokens int, model string) (uint64, error) {
	dollarsPer1KTokens, err := ModelPricingDollarsPer1KTokens(model)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	dollarsPerToken := dollarsPer1KTokens / 1000
	return DollarsToMsats(float64(tokens) * dollarsPerToken), nil
}

// Sha256Hash returns the SHA256 hash in hex of the input hex
func Sha256Hash(hexString string) string {
	bytes, _ := hex.DecodeString(hexString)
	// Create a new SHA256 hash
	h := sha256.New()

	// Write the input hex to the hash
	h.Write(bytes)

	return fmt.Sprintf("%x", h.Sum(nil))
}

// OpenAI Cookbook: https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
func NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return
	}

	var tokensPerMessage, tokensPerName int
	switch model {
	case "gpt-3.5-turbo-0613",
		"gpt-3.5-turbo-16k-0613",
		"gpt-4-0314",
		"gpt-4-32k-0314",
		"gpt-4-0613",
		"gpt-4-32k-0613":
		tokensPerMessage = 3
		tokensPerName = 1
	case "gpt-3.5-turbo-0301":
		tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
		tokensPerName = -1   // if there's a name, the role is omitted
	default:
		if strings.Contains(model, "gpt-3.5-turbo") {
			log.Println("warning: gpt-3.5-turbo may update over time. Returning num tokens assuming gpt-3.5-turbo-0613.")
			return NumTokensFromMessages(messages, "gpt-3.5-turbo-0613")
		} else if strings.Contains(model, "gpt-4") {
			log.Println("warning: gpt-4 may update over time. Returning num tokens assuming gpt-4-0613.")
			return NumTokensFromMessages(messages, "gpt-4-0613")
		} else {
			err = fmt.Errorf("num_tokens_from_messages() is not implemented for model %s. See https://github.com/openai/openai-python/blob/main/chatml.md for information on how messages are converted to tokens.", model)
			log.Println(err)
			return
		}
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		numTokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tokensPerName
		}
	}
	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>
	return numTokens
}
