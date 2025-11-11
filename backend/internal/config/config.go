package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Logging  LoggingConfig
}

// AppConfig contains application-level settings
type AppConfig struct {
	Name        string
	Version     string
	Environment string
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Driver string
	Path   string
}

// AuthConfig contains authentication settings
type AuthConfig struct {
	JWTSecret           string
	JWTExpirationHours  int
	JWTRefreshHours     int
	BcryptCost          int
	SessionTimeout      time.Duration
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level       string
	Development bool
}

var GlobalConfig *Config

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read config file if provided
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Override with environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("STUMPFWORKS")

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	GlobalConfig = &cfg
	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "Stumpf.Works NAS")
	v.SetDefault("app.version", "0.1.0-alpha")
	v.SetDefault("app.environment", "development")

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.readTimeout", "15s")
	v.SetDefault("server.writeTimeout", "15s")
	v.SetDefault("server.idleTimeout", "60s")

	// Database defaults
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.path", "./data/stumpfworks.db")

	// Auth defaults
	v.SetDefault("auth.jwtSecret", generateRandomSecret())
	v.SetDefault("auth.jwtExpirationHours", 24)
	v.SetDefault("auth.jwtRefreshHours", 168) // 7 days
	v.SetDefault("auth.bcryptCost", 12)
	v.SetDefault("auth.sessionTimeout", "24h")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.development", true)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Server.Port)
	}

	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	if c.Database.Path == "" {
		return fmt.Errorf("database path is required")
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// GetServerAddress returns the server address in format "host:port"
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// generateRandomSecret generates a random secret for JWT
// In production, this should be set via environment variable
func generateRandomSecret() string {
	secret := os.Getenv("STUMPFWORKS_AUTH_JWTSECRET")
	if secret != "" {
		return secret
	}
	// Development-only fallback
	return "dev-secret-please-change-in-production"
}
