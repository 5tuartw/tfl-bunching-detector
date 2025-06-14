package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TflKey string
}

func NewConfig() (*Config){
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("FATAL: unable to load .env file: %v.", err)
		os.Exit(1)
	}
	tflKey := os.Getenv("TFL_API_KEY")
	if tflKey == "" {
		log.Fatalf("FATAL: could not find TFL_API_KEY in environment variables")
	}

	cfg := Config{
		TflKey: tflKey,
	}

	return &cfg
}