package utils

import "strings"

// NormalizeEmail converts an email address to lowercase
func NormalizeEmail(email string) string {
	return strings.ToLower(email)
}