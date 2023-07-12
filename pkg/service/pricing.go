package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	models "github.com/kodylow/matador/pkg/models"
	"github.com/kodylow/matador/pkg/utils"
)

func get1000Msats(reqInfo models.RequestInfo) (uint64, error) {
	return 1000, nil
}

func getMsatsChatCompletions(reqInfo models.RequestInfo) (uint64, error) {
	// Get the body of the request as a ChatCompletionRequest
	var body ChatCompletionRequest
	err := json.NewDecoder(bytes.NewReader(reqInfo.Body)).Decode(&body)
	if err != nil {
		log.Println("Error decoding ChatCompletion request body:", err)
		return 0, fmt.Errorf("Error decoding ChatCompletion request body: %w", err)
	}

	// Validate the body
	if err = body.Validate(); err != nil {
		log.Println(err)
		return 0, err
	}

	// Get the number of tokens in the request
	tokens := utils.NumTokensFromMessages(body.Messages, body.Model)

	// tack on tokens for the response, totally naive pricing
	// TODO: best way to rewrite this would be have the LSAT initial purchase be for like a dollar
	// and then subtract off of that for each request when I get the response back
	if tokens < 1000 {
		tokens = 1000
	} else {
		tokens *= 2
	}
	log.Println("Tokens:", tokens)

	// Get the mSats for the number of tokens
	msats, err := utils.TokensToMsats(tokens, body.Model)
	if err != nil {
		log.Println("Error converting tokens to msats:", err)
		return 0, err
	}

	return msats, nil
}

func getMsatsImageGenerations(reqInfo models.RequestInfo) (uint64, error) {
	// Get the body of the request as a ImageGenerationRequest
	var body ImageGenerationRequest
	err := json.NewDecoder(bytes.NewReader(reqInfo.Body)).Decode(&body)
	if err != nil {
		log.Println("Error decoding ImageGeneration request body:", err)
		return 0, fmt.Errorf("Error decoding ImageGeneration request body: %w", err)
	}

	// Validate the body
	if err = body.Validate(); err != nil {
		log.Println(err)
		return 0, err
	}

	var dollars float64
	switch body.Size {
	case "256x256":
		dollars = float64(body.N) * 0.016 // $0.016 per image
	case "512x512":
		dollars = float64(body.N) * 0.018 // $0.018 per image
	default:
		dollars = float64(body.N) * 0.02 // $0.020 per image
	}

	// convert cents to msats
	msats := utils.DollarsToMsats(dollars)

	return msats, nil
}
