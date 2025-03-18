package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

func GenerateAppCode() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback in case of error
		return base64.URLEncoding.EncodeToString([]byte(time.Now().String()))[:8]
	}
	return base64.URLEncoding.EncodeToString(b)[:8]
}
