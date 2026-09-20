package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port          string
	WeatherApiKey string
}

func LoadConfig() (*AppConfig, error) {
	_ = godotenv.Load()

	weatherApiKey, err := getRequiredEnv("WEATHER_API_KEY")
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		Port:          getEnvOrDefault("PORT", "8080"),
		WeatherApiKey: weatherApiKey,
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getRequiredEnv(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("%s é obrigatório", key)
}
