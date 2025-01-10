package logger

import (
	"log"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Beretta350/gochat/pkg/util"
)

var once sync.Once

// Init Single execution that initializes the logger by reading configuration from environment variables.
func Init(environment string) {
	once.Do(func() {
		// Read environment variables
		level := util.GetEnv("LOG_LEVEL", "info")          // Default to "info"
		outputPaths := util.GetEnv("LOG_OUTPUT", "stdout") // Default to "stdout"

		// Parse log level
		var logLevel zapcore.Level
		switch level {
		case "debug":
			logLevel = zap.DebugLevel
		case "info":
			logLevel = zap.InfoLevel
		case "warn":
			logLevel = zap.WarnLevel
		case "error":
			logLevel = zap.ErrorLevel
		default:
			log.Fatalf("Invalid log level: %s", level)
		}

		// Set encoding based on environment
		encoding := "json"
		if environment == "dev" || environment == "local" {
			encoding = "console"
		}

		// Parse output paths (comma-separated)
		outputPathsSlice := strings.Split(outputPaths, ",")

		// Configure logger
		zapConfig := zap.Config{
			Level:            zap.NewAtomicLevelAt(logLevel),
			Development:      environment == "dev" || environment == "local",
			Encoding:         encoding,
			OutputPaths:      outputPathsSlice,
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		}

		// Build logger
		logger, err := zapConfig.Build(zap.AddCallerSkip(1))
		if err != nil {
			log.Fatalf("Failed to initialize logger: %v", err)
		}

		zap.ReplaceGlobals(logger) // Set global logger
		logger.Sugar().Info("Logger configured")
	})
}
