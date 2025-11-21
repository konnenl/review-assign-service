package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort string
	DBUrl     string
}

func LoadConfig() (*Config, error) {
	const pth = "config.LoadConfig"
	dbUrl := getEnv("POSTGRES_URL", "")
	if dbUrl == "" {
		return nil, fmt.Errorf("POSTGRES_URL is not set in .env")
	}
	serverPort := getEnv("SERVER_PORT", ":8080")
	if serverPort[0] != ':' {
		serverPort = ":" + serverPort
	}

	return &Config{
		ServerPort: serverPort,
		DBUrl:     dbUrl,
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
