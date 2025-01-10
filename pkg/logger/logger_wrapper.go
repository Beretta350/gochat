package logger

import "go.uber.org/zap"

// Debug logs a debug message.
func Debug(args ...interface{}) {
	zap.S().Debug(args...)
}

// Debugf logs a debug message formating the message according to the format specifier..
func Debugf(msg string, args ...interface{}) {
	zap.S().Debugf(msg, args...)
}

// Info logs an informational message.
func Info(args ...interface{}) {
	zap.S().Info(args...)
}

// Infof logs an informational message formating the message according to the format specifier..
func Infof(msg string, args ...interface{}) {
	zap.S().Infof(msg, args...)
}

// Warn logs a warning message.
func Warn(args ...interface{}) {
	zap.S().Warn(args...)
}

// Warnf logs a warning message formating the message according to the format specifier..
func Warnf(msg string, args ...interface{}) {
	zap.S().Warnf(msg, args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	zap.S().Error(args...)
}

// Errorf logs an error message formating the message according to the format specifier.
func Errorf(msg string, args ...interface{}) {
	zap.S().Errorf(msg, args...)
}

// Fatal logs a fatal message and exits the application and calls os.Exit.
func Fatal(args ...interface{}) {
	zap.S().Fatal(args...)
}

// Fatalf formats the message according to the format specifier and calls os.Exit.
func Fatalf(msg string, args ...interface{}) {
	zap.S().Fatalf(msg, args...)
}
