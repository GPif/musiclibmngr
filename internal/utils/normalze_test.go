package utils

import "testing"

func TestNormalizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "The quick brown fox", expected: "quick brown fox"},
		{input: "Bob & Alice", expected: "bob and alice"},
		{input: "Feat. Artist", expected: "artist"},
		{input: "l'été", expected: "l ete"},
		{input: "  multiple   spaces  ", expected: "multiple spaces"},
	}

	for _, tt := range tests {
		result := NormalizeString(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeString(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
