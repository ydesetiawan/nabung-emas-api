package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port string
	Env  string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret           string
	JWTExpiry           time.Duration
	RefreshTokenExpiry  time.Duration

	// Storage
	StorageType      string
	StoragePath      string
	AWSRegion        string
	AWSBucket        string
	AWSAccessKey     string
	AWSSecretKey     string

	// CORS
	AllowedOrigins string

	// Email
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	EmailFrom    string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtExpiry, err := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	if err != nil {
		jwtExpiry = 24 * time.Hour
	}

	refreshExpiry, err := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))
	if err != nil {
		refreshExpiry = 168 * time.Hour
	}

	return &Config{
		Port:                getEnv("PORT", "8080"),
		Env:                 getEnv("ENV", "development"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		JWTSecret:           getEnv("JWT_SECRET", ""),
		JWTExpiry:           jwtExpiry,
		RefreshTokenExpiry:  refreshExpiry,
		StorageType:         getEnv("STORAGE_TYPE", "local"),
		StoragePath:         getEnv("STORAGE_PATH", "./uploads"),
		AWSRegion:           getEnv("AWS_REGION", ""),
		AWSBucket:           getEnv("AWS_BUCKET", ""),
		AWSAccessKey:        getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretKey:        getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AllowedOrigins:      getEnv("ALLOWED_ORIGINS", "*"),
		SMTPHost:            getEnv("SMTP_HOST", ""),
		SMTPPort:            getEnv("SMTP_PORT", "587"),
		SMTPUsername:        getEnv("SMTP_USERNAME", ""),
		SMTPPassword:        getEnv("SMTP_PASSWORD", ""),
		EmailFrom:           getEnv("EMAIL_FROM", "noreply@emasgo.com"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
