package utils

import (
	"math/rand"
	"strings"
	"time"
)

// NormalizeString converts a string to lowercase for case-insensitive comparisons
func NormalizeString(s string) string {
	return strings.ToLower(s)
}

// GenerateRandomCode generates a random string of digits of a given length.
func GenerateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = byte(rand.Intn(10)) + '0'
	}
	return string(b)
}