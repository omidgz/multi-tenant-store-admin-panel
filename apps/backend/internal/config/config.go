package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server Configuration
	ServerPort string `env:"SERVER_PORT"`

	// AWS Configuration
	AWSRegion string `env:"AWS_REGION"`

	// PostgreSQL Configuration
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDB       string `env:"POSTGRES_DB"`
	PostgresSSLMode  string `env:"POSTGRES_SSL_MODE"`

	// Cognito Configuration (for JWT validation later)
	CognitoUserPoolID string `env:"COGNITO_USER_POOL_ID"`
	CognitoClientID   string `env:"COGNITO_CLIENT_ID"`
	CognitoJWKSURL    string

	// Tenant Mode
	TenantMode string `env:"TENANT_MODE"` // "header" or "jwt"

	// Application
	Environment string `env:"ENVIRONMENT"` // dev, staging, prod
}

// Load loads configuration from .env file and environment variables
func Load() (*Config, error) {
	// Load .env file for local development (ignored in production)
	_ = godotenv.Load()

	cfg := &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		AWSRegion:         getEnv("AWS_REGION", "us-east-1"),
		PostgresHost:      getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:      getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:      getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword:  getEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:        getEnv("POSTGRES_DB", "store_admin"),
		PostgresSSLMode:   getEnv("POSTGRES_SSL_MODE", "disable"),
		CognitoUserPoolID: getEnv("COGNITO_USER_POOL_ID", ""),
		CognitoClientID:   getEnv("COGNITO_CLIENT_ID", ""),
		TenantMode:        getEnv("TENANT_MODE", "header"),
		Environment:       getEnv("ENVIRONMENT", "dev"),
	}

	// Build Cognito JWKS URL if User Pool ID is provided
	if cfg.CognitoUserPoolID != "" {
		cfg.CognitoJWKSURL = fmt.Sprintf(
			"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
			cfg.AWSRegion,
			cfg.CognitoUserPoolID,
		)
	}

	return cfg, nil
}

// getEnv retrieves environment variable with default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetPostgresConnectionString returns the full connection string for PostgreSQL
func (c *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
		c.PostgresSSLMode,
	)
}
