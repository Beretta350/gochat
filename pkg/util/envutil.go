package util

import (
	"os"
	"strings"
)

// GetEnv gets the value of the environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value)
	}
	return defaultValue
}
