package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateSSEToken generates a cryptographically secure random token
// Returns a 64-character hexadecimal string (256-bit security)
func GenerateSSEToken() (string, error) {
	// Generate 32 bytes of random data (256 bits)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate cryptographically secure random token: %w", err)
	}

	// Encode to hexadecimal string (64 characters)
	return hex.EncodeToString(bytes), nil
}

// ValidateToken validates that a token meets security requirements
// Token must be exactly 64 characters and valid hexadecimal
func ValidateToken(token string) error {
	// Check length
	if len(token) != 64 {
		return fmt.Errorf("invalid token length: expected 64 characters, got %d", len(token))
	}

	// Validate hexadecimal format
	if _, err := hex.DecodeString(token); err != nil {
		return fmt.Errorf("invalid token format: must be valid hexadecimal string")
	}

	return nil
}
