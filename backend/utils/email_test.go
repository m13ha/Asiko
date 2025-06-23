package utils

import "testing"

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "already lowercase",
			email:    "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "uppercase",
			email:    "TEST@EXAMPLE.COM",
			expected: "test@example.com",
		},
		{
			name:     "mixed case",
			email:    "Test@Example.com",
			expected: "test@example.com",
		},
		{
			name:     "empty string",
			email:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeEmail(tt.email)
			if result != tt.expected {
				t.Errorf("NormalizeEmail(%q) = %q, want %q", tt.email, result, tt.expected)
			}
		})
	}
}