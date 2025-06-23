package utils

import "strings"

// NormalizeString converts a string to lowercase for case-insensitive comparisons
func NormalizeString(s string) string {
	return strings.ToLower(s)
}