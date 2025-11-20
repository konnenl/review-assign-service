package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	ServerPort string
	DB_url     string
}

func LoadConfig() (*Config, error) {
	const pth = "config.LoadConfig"
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", pth, err)
	}
	dbUrl := getEnv("DATABASE_URL", "")
	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set in .env")
	}
	serverPort := getEnv("SERVER_PORT", ":8080")
	if serverPort[0] != ':' {
		serverPort = ":" + serverPort
	}

	return &Config{
		ServerPort: serverPort,
		DB_url:     dbUrl,
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
