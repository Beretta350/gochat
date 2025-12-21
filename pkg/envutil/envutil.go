package envutil

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// GetEnv gets the value of the environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value)
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	intValue, err := strconv.Atoi(value)
	if exists && err == nil {
		return intValue
	}
	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	boolValue, err := strconv.ParseBool(value)
	if exists && err == nil {
		return boolValue
	}
	return defaultValue
}

func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
