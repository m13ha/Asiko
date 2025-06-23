package utils

import "testing"

func TestNormalizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "already lowercase",
			input:    "single",
			expected: "single",
		},
		{
			name:     "uppercase",
			input:    "GROUP",
			expected: "group",
		},
		{
			name:     "mixed case",
			input:    "SiNgLe",
			expected: "single",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeString(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}