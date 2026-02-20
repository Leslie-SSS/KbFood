package config

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server    ServerConfig    `envconfig:"SERVER"`
	Database  DatabaseConfig  `envconfig:"DB"`
	Platforms PlatformsConfig `envconfig:"PLATFORM"`
	Log       LogConfig       `envconfig:"LOG"`
	BarkURL   string          `envconfig:"BARK_URL"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port            int           `envconfig:"PORT" mapstructure:"port" default:"8080"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" mapstructure:"read_timeout" default:"30s"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" mapstructure:"write_timeout" default:"30s"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" mapstructure:"shutdown_timeout" default:"30s"`
}

// Address returns the server address
func (c *ServerConfig) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}

// DatabaseConfig holds database configuration for SQLite
type DatabaseConfig struct {
	Path            string        `envconfig:"PATH" mapstructure:"path" default:"file:./data/food.db?mode=rwc"`
	MaxOpenConns    int           `envconfig:"MAX_OPEN_CONNS" mapstructure:"max_open_conns" default:"25"`
	MaxIdleConns    int           `envconfig:"MAX_IDLE_CONNS" mapstructure:"max_idle_conns" default:"5"`
	ConnMaxLifetime time.Duration `envconfig:"MAX_LIFETIME" mapstructure:"max_lifetime" default:"1h"`
	ConnMaxIdleTime time.Duration `envconfig:"MAX_IDLE_TIME" mapstructure:"max_idle_time" default:"30m"`

	// Legacy fields for backwards compatibility (not used for SQLite)
	Host     string `envconfig:"HOST" mapstructure:"host"`
	Port     int    `envconfig:"PORT" mapstructure:"port"`
	User     string `envconfig:"USER" mapstructure:"user"`
	Password string `envconfig:"PASSWORD" mapstructure:"password"`
	Name     string `envconfig:"NAME" mapstructure:"name"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string `envconfig:"LEVEL" mapstructure:"level" default:"info"`
	Format     string `envconfig:"FORMAT" mapstructure:"format" default:"json"`
	TimeFormat string `envconfig:"TIME_FORMAT" mapstructure:"time_format" default:""`
}

// PlatformsConfig holds platform-specific configuration
type PlatformsConfig struct {
	TanTanTang TanTanTangConfig `envconfig:"TANTANTANG"`
	DT         DTConfig         `envconfig:"DT"`
	XiaoCan    XiaoCanConfig    `envconfig:"XIAOCAN"`
}

// TanTanTangConfig holds TanTanTang platform configuration
type TanTanTangConfig struct {
	Token     string `envconfig:"TOKEN" mapstructure:"token" default:""`
	SecretKey string `envconfig:"SECRET_KEY" mapstructure:"secret_key" default:""`
	BaseURL   string `envconfig:"BASE_URL" mapstructure:"base_url" default:""`
}

// DTConfig holds DT platform configuration
type DTConfig struct {
	Token string `envconfig:"TOKEN" mapstructure:"token" default:""`
}

// XiaoCanConfig holds XiaoCan platform configuration
type XiaoCanConfig struct {
	XVayne  string `envconfig:"X_VAYNE" mapstructure:"x_vayne" default:""`
	XTeemo  string `envconfig:"X_TEEMO" mapstructure:"x_teemo" default:""`
	XAshe   string `envconfig:"X_ASHE" mapstructure:"x_ashe" default:""`
	XNami   string `envconfig:"X_NAMI" mapstructure:"x_nami" default:""`
	XSivir  string `envconfig:"X_SIVIR" mapstructure:"x_sivir" default:""`
	UserID  string `envconfig:"USER_ID" mapstructure:"user_id" default:""`
	SilkID  string `envconfig:"SILK_ID" mapstructure:"silk_id" default:""`
}

// Load loads configuration from environment variables and optional config file
func Load(configPath string) (*Config, error) {
	// Try to load config file if provided
	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	} else {
		// Try to find config file in current directory
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./deployments")
		viper.AddConfigPath("/etc/kbfood")

		// Read config file if it exists (don't fail if not found)
		_ = viper.ReadInConfig()
	}

	// Set defaults
	setDefaults()

	// First, unmarshal from config file/viper to get defaults
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: config unmarshal error: %v\n", err)
	}

	// Then load from environment variables using envconfig
	// We use a two-step process:
	// 1. Load env vars into a temporary struct
	// 2. Merge non-empty values into cfg

	type EnvConfig struct {
		ServerPort      int    `envconfig:"SERVER_PORT"`
		DBPath          string `envconfig:"DB_PATH"`
		LogLevel        string `envconfig:"LOG_LEVEL"`
		TantantangToken string `envconfig:"PLATFORM_TANTANTANG_TOKEN"`
		TantantangSK    string `envconfig:"PLATFORM_TANTANTANG_SECRET_KEY"`
		TantantangURL   string `envconfig:"PLATFORM_TANTANTANG_BASE_URL"`
		DTToken         string `envconfig:"PLATFORM_DT_TOKEN"`
		BarkURL         string `envconfig:"BARK_URL"`
	}

	var envCfg EnvConfig
	if err := envconfig.Process("food", &envCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: env config error: %v\n", err)
	}

	// Merge env vars into cfg (only override if env var is set)
	if envCfg.ServerPort != 0 {
		cfg.Server.Port = envCfg.ServerPort
	}
	if envCfg.DBPath != "" {
		cfg.Database.Path = envCfg.DBPath
	}
	if envCfg.LogLevel != "" {
		cfg.Log.Level = envCfg.LogLevel
	}
	if envCfg.TantantangToken != "" {
		cfg.Platforms.TanTanTang.Token = envCfg.TantantangToken
	}
	if envCfg.TantantangSK != "" {
		cfg.Platforms.TanTanTang.SecretKey = envCfg.TantantangSK
	}
	if envCfg.TantantangURL != "" {
		cfg.Platforms.TanTanTang.BaseURL = envCfg.TantantangURL
	}
	if envCfg.DTToken != "" {
		cfg.Platforms.DT.Token = envCfg.DTToken
	}
	if envCfg.BarkURL != "" {
		cfg.BarkURL = envCfg.BarkURL
	}

	// Validate
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	// Database defaults (SQLite)
	viper.SetDefault("database.path", "file:./data/food.db?mode=rwc")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
}

func validate(cfg *Config) error {
	// Validate server port
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	// Validate database path
	if cfg.Database.Path == "" {
		return fmt.Errorf("database path is required")
	}

	return nil
}
