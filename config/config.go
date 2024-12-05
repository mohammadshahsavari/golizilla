package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	Host string
	Port int
	Env  string

	DBUsername string
	DBPassword string
	DBHost     string
	DBName     string
	DBPort     int

	MongoDbUsername    string
	MongoDbPassword    string
	MongoDbHost        string
	MongoDbPort        int
	MongoDbArchivePath string

	EmailSender       string
	EmailSMTPHost     string
	EmailSMTPPort     int
	EmailSMTPUsername string
	EmailSMTPPassword string

	JWTSecretKey string
	JWTExpiresIn time.Duration

	TwoFAExpiresIn        time.Duration
	VerificationExpiresIn time.Duration
}

// LoadConfig loads environment variables from the .env file and returns a Config struct
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Create Config instance from environment variables
	cfg := &Config{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnvAsInt("PORT", 8080),
		Env:  getEnv("ENV", "development"),

		DBUsername: getEnv("DB_USERNAME", "username"),
		DBPassword: getEnv("DB_PASSWORD", "password123"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBName:     getEnv("DB_NAME", "db"),

		MongoDbUsername:    getEnv("MONGODB_USERNAME", "username"),
		MongoDbPassword:    getEnv("MONGODB_PASSWORD", "password123"),
		MongoDbHost:        getEnv("MONGODB_HOST", "localhost"),
		MongoDbPort:        getEnvAsInt("MONGODB_PORT", 27107),
		MongoDbArchivePath: getEnv("MONGODB_ARCHIVE_PATH", "./archive"),

		EmailSender:       getEnv("EMAIL_SENDER", "no-reply@example.com"),
		EmailSMTPHost:     getEnv("EMAIL_SMTP_HOST", "smtp.example.com"),
		EmailSMTPPort:     getEnvAsInt("EMAIL_SMTP_PORT", 587),
		EmailSMTPUsername: getEnv("EMAIL_SMTP_USERNAME", ""),
		EmailSMTPPassword: getEnv("EMAIL_SMTP_PASSWORD", ""),

		JWTSecretKey: getEnv("JWT_SECRET_KEY", "your-default-secret-key"),
		JWTExpiresIn: time.Duration(getEnvAsInt("JWT_EXPIRES_IN", 86400)) * time.Second,

		TwoFAExpiresIn:        time.Duration(getEnvAsInt("2FA_EXPIRES_IN", 600)) * time.Second,
		VerificationExpiresIn: time.Duration(getEnvAsInt("VERIFICATION_EXPIRES_IN", 900)) * time.Second,
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
