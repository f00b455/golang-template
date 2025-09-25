package shared

import (
	"fmt"
	"strings"
)

// Greet returns a greeting message for the given name.
// Returns an error message if the name is empty or only whitespace.
func Greet(name string) string {
	if strings.TrimSpace(name) == "" {
		return "Error: Name cannot be empty"
	}
	return fmt.Sprintf("Hello, %s!", name)
}
