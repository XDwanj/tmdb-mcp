package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				TMDB: TMDBConfig{
					APIKey:    "test_api_key_12345678",
					Language:  "en-US",
					RateLimit: 40,
				},
				Server: ServerConfig{
					Mode: "stdio",
					SSE: SSEConfig{
						Enabled: false,
						Host:    "0.0.0.0",
						Port:    8910,
					},
				},
				Logging: LogConfig{
					Level: "info",
				},
			},
			wantErr: false,
		},
		{
			name: "missing api key",
			config: Config{
				TMDB: TMDBConfig{
					APIKey:    "",
					Language:  "en-US",
					RateLimit: 40,
				},
				Server: ServerConfig{
					Mode: "stdio",
				},
				Logging: LogConfig{
					Level: "info",
				},
			},
			wantErr: true,
			errMsg:  "missing required configuration: TMDB API Key",
		},
		{
			name: "invalid rate limit",
			config: Config{
				TMDB: TMDBConfig{
					APIKey:    "test_api_key",
					Language:  "en-US",
					RateLimit: 0,
				},
				Server: ServerConfig{
					Mode: "stdio",
				},
				Logging: LogConfig{
					Level: "info",
				},
			},
			wantErr: true,
			errMsg:  "invalid rate_limit",
		},
		{
			name: "invalid logging level",
			config: Config{
				TMDB: TMDBConfig{
					APIKey:    "test_api_key",
					Language:  "en-US",
					RateLimit: 40,
				},
				Server: ServerConfig{
					Mode: "stdio",
				},
				Logging: LogConfig{
					Level: "invalid",
				},
			},
			wantErr: true,
			errMsg:  "invalid logging level",
		},
		{
			name: "invalid server mode",
			config: Config{
				TMDB: TMDBConfig{
					APIKey:    "test_api_key",
					Language:  "en-US",
					RateLimit: 40,
				},
				Server: ServerConfig{
					Mode: "invalid",
				},
				Logging: LogConfig{
					Level: "info",
				},
			},
			wantErr: true,
			errMsg:  "invalid server mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLoad_Defaults(t *testing.T) {
	// 使用临时目录避免影响真实配置
	tempDir := t.TempDir()

	// 保存原始 HOME 环境变量
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// 设置临时 HOME
	os.Setenv("HOME", tempDir)

	// 清除可能影响测试的环境变量
	os.Unsetenv("TMDB_API_KEY")
	os.Unsetenv("TMDB_LANGUAGE")
	os.Unsetenv("TMDB_RATE_LIMIT")
	os.Unsetenv("LOGGING_LEVEL")
	os.Unsetenv("SERVER_MODE")

	// 设置必需的 API Key
	os.Setenv("TMDB_API_KEY", "test_key_for_defaults")
	defer os.Unsetenv("TMDB_API_KEY")

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证默认值
	assert.Equal(t, "test_key_for_defaults", cfg.TMDB.APIKey)
	assert.Equal(t, "en-US", cfg.TMDB.Language)
	assert.Equal(t, 40, cfg.TMDB.RateLimit)
	assert.Equal(t, "stdio", cfg.Server.Mode)
	assert.Equal(t, false, cfg.Server.SSE.Enabled)
	assert.Equal(t, "0.0.0.0", cfg.Server.SSE.Host)
	assert.Equal(t, 8910, cfg.Server.SSE.Port)
	assert.Equal(t, "info", cfg.Logging.Level)
}

func TestLoad_EnvironmentVariables(t *testing.T) {
	// 使用临时目录
	tempDir := t.TempDir()

	// 保存原始环境变量
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// 设置临时 HOME
	os.Setenv("HOME", tempDir)

	// 设置环境变量
	testEnvVars := map[string]string{
		"TMDB_API_KEY":       "env_api_key_12345",
		"TMDB_LANGUAGE":      "zh-CN",
		"TMDB_RATE_LIMIT":    "50",
		"LOGGING_LEVEL":      "debug",
		"SERVER_MODE":        "sse",
		"SERVER_SSE_HOST":    "127.0.0.1",
		"SERVER_SSE_PORT":    "9000",
		"SERVER_SSE_ENABLED": "true",
	}

	for k, v := range testEnvVars {
		os.Setenv(k, v)
		defer os.Unsetenv(k)
	}

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证环境变量被正确读取
	assert.Equal(t, "env_api_key_12345", cfg.TMDB.APIKey)
	assert.Equal(t, "zh-CN", cfg.TMDB.Language)
	assert.Equal(t, 50, cfg.TMDB.RateLimit)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "sse", cfg.Server.Mode)
	assert.Equal(t, "127.0.0.1", cfg.Server.SSE.Host)
	assert.Equal(t, 9000, cfg.Server.SSE.Port)
}

func TestLoad_ConfigFile(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".tmdb-mcp")

	// 创建配置目录
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	// 创建配置文件
	configContent := `
tmdb:
  api_key: "file_api_key_abcdef"
  language: "ja-JP"
  rate_limit: 30

server:
  mode: "both"
  sse:
    enabled: true
    host: "localhost"
    port: 8888

logging:
  level: "warn"
`
	configFile := filepath.Join(configDir, "config.yaml")
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// 保存并设置 HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// 清除环境变量，确保从文件读取
	os.Unsetenv("TMDB_API_KEY")
	os.Unsetenv("TMDB_LANGUAGE")
	os.Unsetenv("TMDB_RATE_LIMIT")
	os.Unsetenv("LOGGING_LEVEL")
	os.Unsetenv("SERVER_MODE")

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证从文件读取的配置
	assert.Equal(t, "file_api_key_abcdef", cfg.TMDB.APIKey)
	assert.Equal(t, "ja-JP", cfg.TMDB.Language)
	assert.Equal(t, 30, cfg.TMDB.RateLimit)
	assert.Equal(t, "both", cfg.Server.Mode)
	assert.Equal(t, true, cfg.Server.SSE.Enabled)
	assert.Equal(t, "localhost", cfg.Server.SSE.Host)
	assert.Equal(t, 8888, cfg.Server.SSE.Port)
	assert.Equal(t, "warn", cfg.Logging.Level)
}

func TestLoad_EnvironmentOverridesFile(t *testing.T) {
	// 创建临时目录和配置文件
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".tmdb-mcp")
	err := os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	configContent := `
tmdb:
  api_key: "file_api_key"
  language: "en-US"
  rate_limit: 40
logging:
  level: "info"
server:
  mode: "stdio"
`
	configFile := filepath.Join(configDir, "config.yaml")
	err = os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// 设置 HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	os.Setenv("HOME", tempDir)

	// 设置环境变量覆盖文件值
	os.Setenv("TMDB_API_KEY", "env_overrides_file")
	defer os.Unsetenv("TMDB_API_KEY")
	os.Setenv("TMDB_LANGUAGE", "fr-FR")
	defer os.Unsetenv("TMDB_LANGUAGE")

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// 验证环境变量优先级高于文件
	assert.Equal(t, "env_overrides_file", cfg.TMDB.APIKey)
	assert.Equal(t, "fr-FR", cfg.TMDB.Language)
	// 未被环境变量覆盖的值应该来自文件
	assert.Equal(t, 40, cfg.TMDB.RateLimit)
}

func TestEnsureConfigDir(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".tmdb-mcp")

	// 确认目录不存在
	_, err := os.Stat(configDir)
	assert.True(t, os.IsNotExist(err))

	// 调用 ensureConfigDir
	err = ensureConfigDir(configDir)
	require.NoError(t, err)

	// 验证目录已创建
	info, err := os.Stat(configDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// 验证权限
	assert.Equal(t, os.FileMode(0755), info.Mode().Perm())

	// 再次调用应该成功（幂等性）
	err = ensureConfigDir(configDir)
	require.NoError(t, err)
}

func TestGetConfigDir(t *testing.T) {
	configDir, err := getConfigDir()
	require.NoError(t, err)

	// 验证路径格式
	assert.Contains(t, configDir, ".tmdb-mcp")
	assert.True(t, filepath.IsAbs(configDir))
}
