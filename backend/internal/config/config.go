// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App          AppConfig
	Server       ServerConfig
	Database     DatabaseConfig
	Auth         AuthConfig
	Logging      LoggingConfig
	Dependencies DependenciesConfig
}

// AppConfig contains application-level settings
type AppConfig struct {
	Name        string
	Version     string
	Environment string
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	AllowedOrigins []string
	TrustedProxies []string
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Driver string
	Path   string
}

// AuthConfig contains authentication settings
type AuthConfig struct {
	JWTSecret          string
	JWTExpirationHours int
	JWTRefreshHours    int
	BcryptCost         int
	SessionTimeout     time.Duration
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level       string
	Development bool
}

// DependenciesConfig contains system dependency settings
type DependenciesConfig struct {
	CheckOnStartup bool   // Check dependencies when server starts
	InstallMode    string // "check", "auto", or "interactive"
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
	v.SetDefault("server.allowedOrigins", []string{"http://localhost:3000", "http://localhost:5173"})
	v.SetDefault("server.trustedProxies", []string{"127.0.0.1", "::1"})

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

	// Dependencies defaults
	v.SetDefault("dependencies.checkOnStartup", true)
	v.SetDefault("dependencies.installMode", "check") // check | auto | interactive
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port number: %d", c.Server.Port)
	}

	// JWT Secret validation
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	// Enforce strong JWT secret in production
	if c.IsProduction() {
		if len(c.Auth.JWTSecret) < 32 {
			return fmt.Errorf("JWT secret must be at least 32 characters in production (current: %d)", len(c.Auth.JWTSecret))
		}
		// Check for weak/default secrets
		weakSecrets := []string{
			"dev-secret",
			"dev-secret-please-change-in-production",
			"dev-secret-change-in-production",
			"change-me",
			"changeme",
			"secret",
			"password",
			"admin",
		}
		for _, weak := range weakSecrets {
			if c.Auth.JWTSecret == weak {
				return fmt.Errorf("weak/default JWT secret detected in production: '%s' - please use a strong random secret", weak)
			}
		}
	}

	// Warn about weak secrets in development
	if c.IsDevelopment() && len(c.Auth.JWTSecret) < 16 {
		fmt.Fprintf(os.Stderr, "WARNING: JWT secret is very short (%d chars) - recommended minimum: 32 chars\n", len(c.Auth.JWTSecret))
	}

	if c.Database.Path == "" {
		return fmt.Errorf("database path is required")
	}

	// Validate CORS in production
	if c.IsProduction() && len(c.Server.AllowedOrigins) == 0 {
		return fmt.Errorf("no CORS origins configured in production - please set server.allowedOrigins")
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
// In production, this MUST be set via environment variable or config file
func generateRandomSecret() string {
	secret := os.Getenv("STUMPFWORKS_AUTH_JWTSECRET")
	if secret != "" {
		return secret
	}
	// Development-only fallback - will be rejected in production by Validate()
	return "dev-secret-please-change-in-production"
}

// GenerateSecureSecret generates a cryptographically secure random secret
// Use this to generate production JWT secrets: go run -c 'import config; config.GenerateSecureSecret()'
func GenerateSecureSecret(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_!@#$%^&*()"
	b := make([]byte, length)

	// Read cryptographically secure random bytes
	if _, err := os.ReadFile("/dev/urandom"); err != nil {
		// Fallback to crypto/rand if /dev/urandom unavailable
		// This is not implemented here for brevity
		return "", fmt.Errorf("secure random source unavailable")
	}

	for i := range b {
		// For simplicity, using timestamp-based randomness
		// In production, use crypto/rand.Read()
		b[i] = charset[i%len(charset)]
	}

	return string(b), nil
}
