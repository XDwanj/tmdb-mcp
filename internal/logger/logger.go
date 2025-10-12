package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/XDwanj/tmdb-mcp/internal/config"
)

// InitLogger initializes and returns a configured zap logger
func InitLogger(cfg config.LogConfig) (*zap.Logger, error) {
	// Parse log level from config
	level, err := parseLogLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	// Use development mode for debug level, production mode for others
	var logger *zap.Logger
	if cfg.Level == "debug" {
		// Development mode: console output with color
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("failed to create development logger: %w", err)
		}
	} else {
		// Production mode: JSON output
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(level)
		logger, err = config.Build()
		if err != nil {
			return nil, fmt.Errorf("failed to create production logger: %w", err)
		}
	}

	return logger, nil
}

// parseLogLevel converts a string log level to zapcore.Level
// Supports case-insensitive input (e.g., "DEBUG", "Info", "warn")
func parseLogLevel(level string) (zapcore.Level, error) {
	// Normalize to lowercase for case-insensitive comparison
	normalizedLevel := ""
	for _, r := range level {
		if r >= 'A' && r <= 'Z' {
			normalizedLevel += string(r + 32) // Convert to lowercase
		} else {
			normalizedLevel += string(r)
		}
	}

	switch normalizedLevel {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid log level: %s (must be one of: debug, info, warn, error)", level)
	}
}

// maskSensitiveData is a common function for masking sensitive data
// Shows only first 8 characters followed by "...", or "***" if too short
func maskSensitiveData(data string) string {
	if len(data) > 8 {
		return data[:8] + "..."
	}
	return "***"
}

// MaskAPIKey masks an API key for logging, showing only first 8 characters
// This is a public function that can be used by other packages
func MaskAPIKey(apiKey string) string {
	return maskSensitiveData(apiKey)
}

// MaskToken masks a token for logging, showing only first 8 characters
// This is a public function that can be used by other packages
func MaskToken(token string) string {
	return maskSensitiveData(token)
}
