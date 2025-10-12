package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/XDwanj/tmdb-mcp/internal/config"
)

func TestInitLogger_DevelopmentMode(t *testing.T) {
	cfg := config.LogConfig{
		Level: "debug",
	}

	logger, err := InitLogger(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// Cleanup
	if logger != nil {
		logger.Sync()
	}
}

func TestInitLogger_ProductionMode(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"Info level", "info"},
		{"Warn level", "warn"},
		{"Error level", "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.LogConfig{
				Level: tt.level,
			}

			logger, err := InitLogger(cfg)
			assert.NoError(t, err)
			assert.NotNil(t, logger)

			// Cleanup
			if logger != nil {
				logger.Sync()
			}
		})
	}
}

func TestInitLogger_InvalidLevel(t *testing.T) {
	cfg := config.LogConfig{
		Level: "invalid",
	}

	logger, err := InitLogger(cfg)
	assert.Error(t, err)
	assert.Nil(t, logger)
}

func TestParseLogLevel_ValidLevels(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected zapcore.Level
	}{
		{"Debug level", "debug", zapcore.DebugLevel},
		{"Info level", "info", zapcore.InfoLevel},
		{"Warn level", "warn", zapcore.WarnLevel},
		{"Error level", "error", zapcore.ErrorLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := parseLogLevel(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, level)
		})
	}
}

func TestParseLogLevel_InvalidLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Invalid string", "invalid"},
		{"Random uppercase", "INVALID"},
		{"Random mixed case", "Invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := parseLogLevel(tt.input)
			assert.Error(t, err)
			// Should default to InfoLevel even on error
			assert.Equal(t, zapcore.InfoLevel, level)
		})
	}
}

// TestParseLogLevel_CaseInsensitive tests case-insensitive log level parsing
func TestParseLogLevel_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected zapcore.Level
	}{
		{"Uppercase DEBUG", "DEBUG", zapcore.DebugLevel},
		{"Uppercase INFO", "INFO", zapcore.InfoLevel},
		{"Uppercase WARN", "WARN", zapcore.WarnLevel},
		{"Uppercase ERROR", "ERROR", zapcore.ErrorLevel},
		{"Mixed case Debug", "Debug", zapcore.DebugLevel},
		{"Mixed case Info", "Info", zapcore.InfoLevel},
		{"Mixed case Warn", "Warn", zapcore.WarnLevel},
		{"Mixed case Error", "Error", zapcore.ErrorLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := parseLogLevel(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, level)
		})
	}
}

func TestMaskAPIKey(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal API key (> 8 chars)",
			input:    "abcdefghijklmnop",
			expected: "abcdefgh...",
		},
		{
			name:     "Exactly 8 chars",
			input:    "12345678",
			expected: "***",
		},
		{
			name:     "Short API key (< 8 chars)",
			input:    "short",
			expected: "***",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "***",
		},
		{
			name:     "Very long API key",
			input:    "1234567890abcdefghijklmnopqrstuvwxyz",
			expected: "12345678...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskAPIKey(tt.input)
			assert.Equal(t, tt.expected, result)

			// Security check: ensure masked value doesn't contain full original
			if len(tt.input) > 8 {
				assert.NotEqual(t, tt.input, result)
				assert.NotContains(t, result, tt.input[8:])
			}
		})
	}
}

func TestMaskToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal token (> 8 chars)",
			input:    "token123456789",
			expected: "token123...",
		},
		{
			name:     "Exactly 8 chars",
			input:    "token123",
			expected: "***",
		},
		{
			name:     "Short token (< 8 chars)",
			input:    "token",
			expected: "***",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "***",
		},
		{
			name:     "Very long token",
			input:    "abcdefghijklmnopqrstuvwxyz0123456789",
			expected: "abcdefgh...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskToken(tt.input)
			assert.Equal(t, tt.expected, result)

			// Security check: ensure masked value doesn't contain full original
			if len(tt.input) > 8 {
				assert.NotEqual(t, tt.input, result)
				assert.NotContains(t, result, tt.input[8:])
			}
		})
	}
}

// TestInitLogger_MultipleCalls tests that multiple logger initializations work
func TestInitLogger_MultipleCalls(t *testing.T) {
	cfg := config.LogConfig{
		Level: "info",
	}

	logger1, err1 := InitLogger(cfg)
	assert.NoError(t, err1)
	assert.NotNil(t, logger1)

	logger2, err2 := InitLogger(cfg)
	assert.NoError(t, err2)
	assert.NotNil(t, logger2)

	// Cleanup
	if logger1 != nil {
		logger1.Sync()
	}
	if logger2 != nil {
		logger2.Sync()
	}
}

// TestMaskAPIKey_DoesNotLeakSensitiveData verifies masking security
func TestMaskAPIKey_DoesNotLeakSensitiveData(t *testing.T) {
	sensitiveKey := "secretapikey12345678"
	masked := MaskAPIKey(sensitiveKey)

	// Ensure masked value is not the same as original
	assert.NotEqual(t, sensitiveKey, masked)

	// Ensure masked value doesn't contain the sensitive part
	sensitivePart := sensitiveKey[8:]
	assert.NotContains(t, masked, sensitivePart)

	// Ensure it only shows first 8 chars + "..."
	assert.Equal(t, "secretap...", masked)
}

// TestMaskToken_DoesNotLeakSensitiveData verifies token masking security
func TestMaskToken_DoesNotLeakSensitiveData(t *testing.T) {
	sensitiveToken := "mysecrettoken123456"
	masked := MaskToken(sensitiveToken)

	// Ensure masked value is not the same as original
	assert.NotEqual(t, sensitiveToken, masked)

	// Ensure masked value doesn't contain the sensitive part
	sensitivePart := sensitiveToken[8:]
	assert.NotContains(t, masked, sensitivePart)

	// Ensure it only shows first 8 chars + "..."
	assert.Equal(t, "mysecret...", masked)
}
