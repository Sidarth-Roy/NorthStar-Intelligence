package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
	AppEnv      string
}

func LoadConfig() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "host=localhost user=user password=password dbname=northwind port=5432 sslmode=disable"),
		Port:        getEnv("PORT", "8081"),
		AppEnv:      getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}