package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config is the root configuration structure
type Config struct {
	TMDB    TMDBConfig   `mapstructure:"tmdb" json:"tmdb"`
	Server  ServerConfig `mapstructure:"server" json:"server"`
	Logging LogConfig    `mapstructure:"logging" json:"logging"`

	// TokenGenerated indicates if the SSE token was auto-generated
	// This is not persisted to config file
	TokenGenerated bool `mapstructure:"-" json:"-"`
}

// TMDBConfig contains TMDB API configuration
type TMDBConfig struct {
	APIKey    string `mapstructure:"api_key" json:"api_key"`
	Language  string `mapstructure:"language" json:"language"`
	RateLimit int    `mapstructure:"rate_limit" json:"rate_limit"`
}

// ServerConfig contains server configuration
type ServerConfig struct {
	Mode string    `mapstructure:"mode" json:"mode"`
	SSE  SSEConfig `mapstructure:"sse" json:"sse"`
}

// SSEConfig contains SSE server configuration
type SSEConfig struct {
	// Enabled bool   `mapstructure:"enabled" json:"enabled"`
	Host  string `mapstructure:"host" json:"host"`
	Port  int    `mapstructure:"port" json:"port"`
	Token string `mapstructure:"token" json:"token"`
}

// LogConfig contains logging configuration
type LogConfig struct {
	Level string `mapstructure:"level" json:"level"`
}

// Load loads configuration from multiple sources with priority: CLI > ENV > File
func Load() (*Config, error) {
	v := viper.New()

	// 设置默认值
	setDefaults(v)

	// 配置文件路径
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	// 如果配置目录不存在，自动创建
	if err := ensureConfigDir(configDir); err != nil {
		return nil, fmt.Errorf("failed to ensure config directory: %w", err)
	}

	// 设置配置文件路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)

	// 读取配置文件（允许不存在）
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// 配置文件不存在是可以接受的，使用默认值和环境变量
	}

	// 绑定环境变量
	v.SetEnvPrefix("TMDB")
	v.AutomaticEnv()

	// 手动绑定所有配置键到环境变量
	bindEnvVars(v)

	// 解析配置到结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 处理 SSE Token（如果 SSE 模式启用）
	if cfg.Server.Mode == "sse" || cfg.Server.Mode == "both" {
		if err := handleSSEToken(&cfg, v, configDir); err != nil {
			return nil, fmt.Errorf("failed to handle SSE token: %w", err)
		}
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// 检查必需配置：TMDB API Key
	if c.TMDB.APIKey == "" {
		return fmt.Errorf("missing required configuration: TMDB API Key. Please set TMDB_API_KEY environment variable or add it to ~/.tmdb-mcp/config.yaml")
	}

	// 检查 rate limit 有效性
	if c.TMDB.RateLimit <= 0 {
		return fmt.Errorf("invalid rate_limit: must be greater than 0")
	}

	// 检查日志级别有效性
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("invalid logging level: %s (must be one of: debug, info, warn, error)", c.Logging.Level)
	}

	// 检查 server mode 有效性
	validModes := map[string]bool{
		"stdio": true,
		"sse":   true,
		"both":  true,
	}
	if !validModes[c.Server.Mode] {
		return fmt.Errorf("invalid server mode: %s (must be one of: stdio, sse, both)", c.Server.Mode)
	}

	return nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// TMDB defaults
	v.SetDefault("tmdb.language", "en-US")
	v.SetDefault("tmdb.rate_limit", 40)

	// Server defaults
	v.SetDefault("server.mode", "both")
	v.SetDefault("server.sse.enabled", false)
	v.SetDefault("server.sse.host", "0.0.0.0")
	v.SetDefault("server.sse.port", 8910)

	// Logging defaults
	v.SetDefault("logging.level", "info")
}

// bindEnvVars binds all configuration keys to environment variables
func bindEnvVars(v *viper.Viper) {
	// TMDB
	v.BindEnv("tmdb.api_key", "TMDB_API_KEY")
	v.BindEnv("tmdb.language", "TMDB_LANGUAGE")
	v.BindEnv("tmdb.rate_limit", "TMDB_RATE_LIMIT")

	// Server
	v.BindEnv("server.mode", "SERVER_MODE")
	v.BindEnv("server.sse.enabled", "SERVER_SSE_ENABLED")
	v.BindEnv("server.sse.host", "SERVER_SSE_HOST")
	v.BindEnv("server.sse.port", "SERVER_SSE_PORT")
	v.BindEnv("server.sse.token", "SSE_TOKEN")

	// Logging
	v.BindEnv("logging.level", "LOGGING_LEVEL")
}

// getConfigDir returns the configuration directory path
func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, ".tmdb-mcp"), nil
}

// ensureConfigDir ensures the configuration directory exists
func ensureConfigDir(configDir string) error {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	}
	return nil
}

// handleSSEToken handles SSE token loading, generation, and validation
func handleSSEToken(cfg *Config, v *viper.Viper, configDir string) error {
	token := cfg.Server.SSE.Token

	// If token is empty, generate a new one
	if token == "" {
		newToken, err := GenerateSSEToken()
		if err != nil {
			return fmt.Errorf("failed to generate SSE token: %w", err)
		}
		token = newToken
		cfg.Server.SSE.Token = token
		cfg.TokenGenerated = true

		// Save the generated token to config file
		if err := SaveTokenToConfig(v, configDir, token); err != nil {
			return fmt.Errorf("failed to save generated token to config: %w", err)
		}
	}

	// Validate token format
	if err := ValidateToken(token); err != nil {
		return fmt.Errorf("invalid SSE token: %w", err)
	}

	return nil
}

// SaveTokenToConfig saves the SSE token to the configuration file
func SaveTokenToConfig(v *viper.Viper, configDir, token string) error {
	// Set the token in viper
	v.Set("server.sse.token", token)

	// Determine config file path
	configFile := v.ConfigFileUsed()
	if configFile == "" {
		configFile = filepath.Join(configDir, "config.yaml")
	}

	// Write config to file
	if err := v.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Set file permissions to 0600 (owner read/write only)
	if err := os.Chmod(configFile, 0600); err != nil {
		return fmt.Errorf("failed to set config file permissions: %w", err)
	}

	return nil
}
