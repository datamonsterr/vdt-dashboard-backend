package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Environment    string
	Port           string
	DatabaseURL    string
	DatabaseHost   string
	DatabasePort   string
	DatabaseUser   string
	DatabasePass   string
	DatabaseName   string
	LogLevel       string
	AllowOrigins   []string
	ClerkSecretKey string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		DatabaseHost:   getEnv("DB_HOST", "localhost"),
		DatabasePort:   getEnv("DB_PORT", "5432"),
		DatabaseUser:   getEnv("DB_USER", "postgres"),
		DatabasePass:   getEnv("DB_PASSWORD", "postgres"),
		DatabaseName:   getEnv("DB_NAME", "vdt_dashboard"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		ClerkSecretKey: getEnv("CLERK_SECRET_KEY", ""),
		AllowOrigins: []string{
			getEnv("FRONTEND_URL", "http://localhost:3000"),
			getEnv("STORYBOOK_URL", "http://localhost:6006"),
		},
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getEnvAsBool gets an environment variable as boolean with a fallback value
func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return fallback
}
