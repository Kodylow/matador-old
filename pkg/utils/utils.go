package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const BTCPRICE float64 = 25000 // $25,000.00 per bitcoin

func CentsToMsats(cents uint64) uint64 {
    btc := float64(cents) / (BTCPRICE * 100) // Convert cents to BTC
    msats := btc * 100000000000 // Convert BTC to msats
    return uint64(msats)
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
