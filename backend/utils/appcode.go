package utils

import (
	"crypto/rand"
	"encoding/base32"
	"strings"

	mathrand "math/rand"
	"time"
)

func GenerateAppCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)

	_, err := rand.Read(b) // Use crypto/rand for secure random bytes
	if err != nil {
		// Fallback in case of error
		result := make([]byte, 8)
		seed := time.Now().UnixNano()
		r := mathrand.NewSource(seed)
		rng := mathrand.New(r)

		for i := range result {
			result[i] = charset[rng.Intn(len(charset))]
		}
		return string(result)
	}

	// Convert random bytes to uppercase alphanumeric
	result := make([]byte, 8)
	for i := range result {
		// Use modulo to map byte to charset index
		result[i] = charset[int(b[i])%len(charset)]
	}

	return string(result)
}

// GenerateBookingCode creates a secure, unique booking code
func GenerateBookingCode() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		panic("unable to generate booking code")
	}
	code := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	return strings.ToUpper(code)
}
