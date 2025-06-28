package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	Environment string
	JWTSecret   string
}

func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=rural_health_db port=5432 sslmode=disable"),
		Port:        getEnv("PORT", "3000"),
		Environment: getEnv("ENVIRONMENT", "development"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
