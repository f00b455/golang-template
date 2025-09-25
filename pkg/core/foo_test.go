package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFooProcess(t *testing.T) {
	tests := []struct {
		name     string
		config   FooConfig
		input    string
		expected string
	}{
		{
			name:     "with prefix and suffix",
			config:   FooConfig{Prefix: "✨", Suffix: "✨"},
			input:    "Hello, World!",
			expected: "✨Hello, World!✨",
		},
		{
			name:     "with prefix only",
			config:   FooConfig{Prefix: ">> ", Suffix: ""},
			input:    "test",
			expected: ">> test",
		},
		{
			name:     "empty config",
			config:   FooConfig{},
			input:    "test",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FooProcess(tt.config, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFooGreet(t *testing.T) {
	config := FooConfig{Prefix: "✨", Suffix: "✨"}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid name",
			input:    "Alice",
			expected: "✨Hello, Alice!✨",
		},
		{
			name:     "empty name",
			input:    "",
			expected: "✨Error: Name cannot be empty✨",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FooGreet(config, tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFooProcessor(t *testing.T) {
	config := FooConfig{Prefix: "[", Suffix: "]"}
	processor := NewFooProcessor(config)

	t.Run("Process", func(t *testing.T) {
		result := processor.Process("test")
		assert.Equal(t, "[test]", result)
	})

	t.Run("GreetWithFoo", func(t *testing.T) {
		result := processor.GreetWithFoo("Bob")
		assert.Equal(t, "[Hello, Bob!]", result)
	})
}

func TestFooTransform(t *testing.T) {
	input := []string{"a", "b", "c"}
	transformer := func(s string) string {
		return s + s
	}

	result := FooTransform(input, transformer)
	expected := []string{"aa", "bb", "cc"}

	assert.Equal(t, expected, result)
}

func TestFooFilter(t *testing.T) {
	t.Run("filter strings", func(t *testing.T) {
		input := []string{"apple", "banana", "apricot", "cherry"}
		predicate := func(s string) bool {
			return s[0] == 'a'
		}

		result := FooFilter(input, predicate)
		expected := []string{"apple", "apricot"}

		assert.Equal(t, expected, result)
	})

	t.Run("filter integers", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		predicate := func(n int) bool {
			return n%2 == 0
		}

		result := FooFilter(input, predicate)
		expected := []int{2, 4, 6}

		assert.Equal(t, expected, result)
	})
}
