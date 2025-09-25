package shared

import (
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "standard date",
			input:    time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			expected: "2023-12-25",
		},
		{
			name:     "leap year",
			input:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			expected: "2024-02-29",
		},
		{
			name:     "single digit month and day",
			input:    time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC),
			expected: "2023-01-05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDate(tt.input)
			if result != tt.expected {
				t.Errorf("FormatDate(%v) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
