package shared

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid name",
			input:    "World",
			expected: "Hello, World!",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "Error: Name cannot be empty",
		},
		{
			name:     "whitespace only",
			input:    "   ",
			expected: "Error: Name cannot be empty",
		},
		{
			name:     "name with spaces",
			input:    "John Doe",
			expected: "Hello, John Doe!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Greet(tt.input)
			if result != tt.expected {
				t.Errorf("Greet(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
