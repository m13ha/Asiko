package utils

import (
	"crypto/rand"
	"encoding/base32"
	"strings"
)

func GenerateCode(codeType string) string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		panic("unable to generate booking code")
	}
	code := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	return codeType + strings.ToUpper(code)
}
