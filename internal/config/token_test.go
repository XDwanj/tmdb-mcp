package config

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateSSEToken tests the SSE token generation functionality
func TestGenerateSSEToken(t *testing.T) {
	tests := []struct {
		name    string
		wantLen int
		wantErr bool
	}{
		{
			name:    "generates valid 64-character token",
			wantLen: 64,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateSSEToken()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, token, tt.wantLen, "token should be 64 characters")

			// 验证十六进制格式
			_, err = hex.DecodeString(token)
			assert.NoError(t, err, "token should be valid hexadecimal")
		})
	}
}

// TestGenerateSSEToken_Randomness tests that generated tokens are unique
func TestGenerateSSEToken_Randomness(t *testing.T) {
	// Generate two tokens
	token1, err := GenerateSSEToken()
	require.NoError(t, err)

	token2, err := GenerateSSEToken()
	require.NoError(t, err)

	// Tokens should be different (randomness check)
	assert.NotEqual(t, token1, token2, "generated tokens should be different")
}

// TestValidateToken tests the token validation functionality
func TestValidateToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid 64-character hex token",
			token:   "a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890",
			wantErr: false,
		},
		{
			name:    "invalid length - too short",
			token:   "a1b2c3d4e5f67890",
			wantErr: true,
			errMsg:  "invalid token length",
		},
		{
			name:    "invalid length - too long",
			token:   "a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890extra",
			wantErr: true,
			errMsg:  "invalid token length",
		},
		{
			name:    "invalid characters - non-hex",
			token:   "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
			wantErr: true,
			errMsg:  "invalid token format",
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
			errMsg:  "invalid token length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGenerateAndValidateToken tests integration between generation and validation
func TestGenerateAndValidateToken(t *testing.T) {
	// Generate a token
	token, err := GenerateSSEToken()
	require.NoError(t, err)

	// Validate the generated token
	err = ValidateToken(token)
	assert.NoError(t, err, "generated token should pass validation")
}
