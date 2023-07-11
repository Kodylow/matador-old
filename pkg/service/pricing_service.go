package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	models "github.com/kodylow/actually_openai/pkg/models"
	"github.com/kodylow/actually_openai/pkg/utils"
)

func getMsats(reqInfo models.RequestInfo) (uint64, error) {
  return 1000, nil
}


func getMsatsImageGenerations(reqInfo models.RequestInfo) (uint64, error) {
    // Get the body of the request as a ImageGenerationBody
    var body ImageGenerationBody
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

    var cents uint64
    switch body.Size {
    case "256x256":
        cents = body.N * 16 // $0.016 per image
    case "512x512":
        cents = body.N * 18 // $0.018 per image
    default:
        cents = body.N * 20 // $0.020 per image
    }

    // convert cents to msats
    msats := utils.CentsToMsats(cents)

    return msats, nil
}