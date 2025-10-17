package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFirstStartup_Integration tests the first startup scenario with token auto-generation
func TestFirstStartup_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temporary config directory
	tmpDir := t.TempDir()

	// Save original environment variables
	originalHome := os.Getenv("HOME")
	originalAPIKey := os.Getenv("TMDB_API_KEY")
	originalServerMode := os.Getenv("SERVER_MODE")
	defer func() {
		os.Setenv("HOME", originalHome)
		if originalAPIKey == "" {
			os.Unsetenv("TMDB_API_KEY")
		} else {
			os.Setenv("TMDB_API_KEY", originalAPIKey)
		}
		if originalServerMode == "" {
			os.Unsetenv("SERVER_MODE")
		} else {
			os.Setenv("SERVER_MODE", originalServerMode)
		}
	}()

	// Set environment for test (no config file, using env vars)
	os.Setenv("HOME", tmpDir)
	os.Setenv("TMDB_API_KEY", "test_api_key_for_integration")
	os.Setenv("SERVER_MODE", "sse")

	// Execute: Load config (should trigger token generation)
	cfg, err := Load()
	require.NoError(t, err)

	// Verify: Token was generated
	assert.NotEmpty(t, cfg.Server.SSE.Token, "token should be generated")
	assert.Len(t, cfg.Server.SSE.Token, 64, "token should be 64 characters")
	assert.True(t, cfg.TokenGenerated, "TokenGenerated flag should be true")

	// Verify: Token was saved to config file
	configFile := filepath.Join(tmpDir, ".tmdb-mcp", "config.yaml")
	_, err = os.Stat(configFile)
	require.NoError(t, err, "config file should be created")

	v := viper.New()
	v.SetConfigFile(configFile)
	err = v.ReadInConfig()
	require.NoError(t, err)

	savedToken := v.GetString("server.sse.token")
	assert.Equal(t, cfg.Server.SSE.Token, savedToken, "token should be saved to config file")

	// Verify: Config file permissions (Linux/macOS only)
	if runtime.GOOS != "windows" {
		info, err := os.Stat(configFile)
		require.NoError(t, err)
		perm := info.Mode().Perm()
		assert.Equal(t, os.FileMode(0600), perm, "config file should have 0600 permissions")
	}
}

// TestEnvironmentVariable_Integration tests token loading from environment variable
func TestEnvironmentVariable_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temporary config directory
	tmpDir := t.TempDir()

	// Set environment variable with test token
	testToken := "a1b2c3d4e5f67890abcdef1234567890abcdef1234567890abcdef1234567890"
	os.Setenv("SSE_TOKEN", testToken)
	defer os.Unsetenv("SSE_TOKEN")

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpDir)

	// Create minimal config file
	configDir := filepath.Join(tmpDir, ".tmdb-mcp")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configFile := filepath.Join(configDir, "config.yaml")
	v := viper.New()
	v.SetConfigFile(configFile)
	v.Set("tmdb.api_key", "test_api_key")
	v.Set("server.mode", "sse")
	require.NoError(t, v.WriteConfig())

	// Execute: Load config
	cfg, err := Load()
	require.NoError(t, err)

	// Verify: Token loaded from environment variable
	assert.Equal(t, testToken, cfg.Server.SSE.Token, "token should be loaded from environment variable")
	assert.False(t, cfg.TokenGenerated, "TokenGenerated should be false (loaded from env)")
}

// TestConfigFile_Integration tests token loading from config file
func TestConfigFile_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temporary config directory
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".tmdb-mcp")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	// Create config file with token
	testToken := "b2c3d4e5f67890ababcdef1234567890abcdef1234567890abcdef1234567890"
	configFile := filepath.Join(configDir, "config.yaml")

	v := viper.New()
	v.SetConfigFile(configFile)
	v.Set("tmdb.api_key", "test_api_key")
	v.Set("server.mode", "sse")
	v.Set("server.sse.token", testToken)
	require.NoError(t, v.WriteConfig())

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpDir)

	// Execute: Load config
	cfg, err := Load()
	require.NoError(t, err)

	// Verify: Token loaded from config file
	assert.Equal(t, testToken, cfg.Server.SSE.Token, "token should be loaded from config file")
	assert.False(t, cfg.TokenGenerated, "TokenGenerated should be false (loaded from file)")
}

// TestTokenPriority_Integration tests the priority: ENV > FILE > Generated
func TestTokenPriority_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temporary config directory
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".tmdb-mcp")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	// Create config file with one token
	fileToken := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	configFile := filepath.Join(configDir, "config.yaml")

	v := viper.New()
	v.SetConfigFile(configFile)
	v.Set("tmdb.api_key", "test_api_key")
	v.Set("server.mode", "sse")
	v.Set("server.sse.token", fileToken)
	require.NoError(t, v.WriteConfig())

	// Set environment variable with different token
	envToken := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	os.Setenv("SSE_TOKEN", envToken)
	defer os.Unsetenv("SSE_TOKEN")

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpDir)

	// Execute: Load config
	cfg, err := Load()
	require.NoError(t, err)

	// Verify: Environment variable takes priority
	assert.Equal(t, envToken, cfg.Server.SSE.Token, "environment variable should take priority over config file")
	assert.False(t, cfg.TokenGenerated, "TokenGenerated should be false")
}

// TestSSEDisabled_Integration tests that token is not required when SSE is disabled
func TestSSEDisabled_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Setup: Create temporary config directory
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".tmdb-mcp")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	// Create config file for stdio mode (no token)
	configFile := filepath.Join(configDir, "config.yaml")

	v := viper.New()
	v.SetConfigFile(configFile)
	v.Set("tmdb.api_key", "test_api_key")
	v.Set("server.mode", "stdio")
	require.NoError(t, v.WriteConfig())

	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tmpDir)

	// Execute: Load config
	cfg, err := Load()
	require.NoError(t, err)

	// Verify: Token should be empty (not required for stdio mode)
	assert.Empty(t, cfg.Server.SSE.Token, "token should not be generated for stdio mode")
	assert.False(t, cfg.TokenGenerated, "TokenGenerated should be false")
}
