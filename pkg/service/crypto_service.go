package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

// GetRequestHash returns the SHA256 hash of the request's relevant fields
func GetRequestHash(r *http.Request) string {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	
	// Create a new SHA256 hash
	h := sha256.New()

	// Write the relevant parts of the request to the hash
	h.Write([]byte(r.Method))
	h.Write([]byte(r.URL.String()))
	h.Write(bodyBytes)

	// Reset the request body to its original state
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return fmt.Sprintf("%x", h.Sum(nil))
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
