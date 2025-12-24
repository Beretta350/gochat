package logger_test

import (
	"os"
	"testing"

	"github.com/Beretta350/gochat/pkg/logger"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	// Save original env values to restore later
	originalLogLevel := os.Getenv("LOG_LEVEL")
	originalLogOutput := os.Getenv("LOG_OUTPUT")
	defer func() {
		os.Setenv("LOG_LEVEL", originalLogLevel)
		os.Setenv("LOG_OUTPUT", originalLogOutput)
	}()

	// Test different log levels and environments
	testCases := []struct {
		name        string
		environment string
		logLevel    string
	}{
		{
			name:        "Development environment with debug level",
			environment: "dev",
			logLevel:    "debug",
		},
		{
			name:        "Production environment with error level",
			environment: "prod",
			logLevel:    "error",
		},
		{
			name:        "Default info level",
			environment: "prod",
			logLevel:    "info",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("LOG_LEVEL", tc.logLevel)
			os.Setenv("LOG_OUTPUT", "stdout")

			// Initialize logger - in a real test, we'd normally initialize
			// just once, but since we're checking different configurations,
			// we'll rely on the global zap logger
			logger.Init(tc.environment)

			// Simple verification that initialization doesn't panic
			// For deeper testing, we'd need to export more internals or use mocks

			// Perform a log operation to verify the logger works
			zap.L().Sugar().Info("Test log message")
		})
	}
}

func TestLoggerSingleton(t *testing.T) {
	// First call
	logger.Init("dev")

	// Get global zap logger after first init
	firstLogger := zap.L()

	// Call again with different parameters
	logger.Init("prod")

	// Get global zap logger after second init
	secondLogger := zap.L()

	// If singleton pattern works, both should be the same instance
	if firstLogger != secondLogger {
		t.Error("Logger was reinitialized, singleton pattern not working")
	}
}

// TestLoggerWrappers tests basic functionality of wrapper functions
func TestLoggerWrappers(t *testing.T) {
	// These tests just verify the wrappers don't panic
	// In a real test, we'd mock the zap logger to verify the right methods are called

	// Initialize logger first
	logger.Init("dev")

	// Test wrapper functions
	logger.Debug("Debug message")
	logger.Debugf("Debug message with %s", "formatting")

	logger.Info("Info message")
	logger.Infof("Info message with %s", "formatting")

	logger.Warn("Warning message")
	logger.Warnf("Warning message with %s", "formatting")

	logger.Error("Error message")
	logger.Errorf("Error message with %s", "formatting")

	// Don't test Fatal functions as they would exit the program
}
