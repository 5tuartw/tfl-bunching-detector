package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TflKey string
}

func NewConfig() (*Config, error){
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("could not load .env file: %w", err)
	}
	tflKey := os.Getenv("TFL_API_KEY")
	if tflKey == "" {
		return nil, fmt.Errorf("could not find TFL_API_KEY in environment variables")
	}

	cfg := Config{
		TflKey: tflKey,
	}

	return &cfg, nil
}