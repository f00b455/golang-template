package core

import "github.com/f00b455/golang-template/pkg/shared"

// FooConfig holds configuration for foo processing.
type FooConfig struct {
	Prefix string
	Suffix string
}

// FooProcess applies prefix and suffix to input string.
func FooProcess(config FooConfig, input string) string {
	return config.Prefix + input + config.Suffix
}

// FooGreet creates a greeting with foo processing.
func FooGreet(config FooConfig, name string) string {
	greeting := shared.Greet(name)
	return FooProcess(config, greeting)
}

// FooProcessor holds configuration and provides processing methods.
type FooProcessor struct {
	config FooConfig
}

// NewFooProcessor creates a new FooProcessor with the given config.
func NewFooProcessor(config FooConfig) *FooProcessor {
	return &FooProcessor{config: config}
}

// Process applies the configuration to the input string.
func (fp *FooProcessor) Process(input string) string {
	return FooProcess(fp.config, input)
}

// GreetWithFoo creates a greeting with foo processing.
func (fp *FooProcessor) GreetWithFoo(name string) string {
	return FooGreet(fp.config, name)
}

// FooTransform applies a transformer function to all items in the slice.
func FooTransform(data []string, transformer func(string) string) []string {
	result := make([]string, len(data))
	for i, item := range data {
		result[i] = transformer(item)
	}
	return result
}

// FooFilter filters items based on the predicate function.
func FooFilter[T any](items []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}
