package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	Host       string
	DBUsername string
	DBPassword string
	DBHost     string
	DBName     string
	Env        string
	Port       int
	DBPort     int

	EmailSender   string
	EmailSMTPHost string
	EmailSMTPPort int
	EmailUsername string
	EmailPassword string
}

// LoadConfig loads environment variables from the .env file and returns a Config struct
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Create Config instance from environment variables
	cfg := &Config{
		Host:       getEnv("HOST", "0.0.0.0"),
		Port:       getEnvAsInt("PORT", 8080),
		DBUsername: getEnv("DB_USERNAME", "username"),
		DBPassword: getEnv("DB_PASSWORD", "password123"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBName:     getEnv("DB_NAME", "db"),
		Env:        getEnv("ENV", "development"),

		EmailSender:   getEnv("EMAIL_SENDER", "no-reply@example.com"),
		EmailSMTPHost: getEnv("EMAIL_SMTP_HOST", "smtp.example.com"),
		EmailSMTPPort: getEnvAsInt("EMAIL_SMTP_PORT", 587),
		EmailUsername: getEnv("EMAIL_USERNAME", ""),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
	}

	return cfg, nil
}

// Helper function to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to read an environment variable as an integer or return a default value
func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
