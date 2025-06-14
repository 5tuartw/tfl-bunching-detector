package main

import (
	"log"

	"github.com/5tuartw/tfl-bunching-detector/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		if cfg.TflKey == "" {
			log.Printf("ERROR: TFL_API_KEY not loaded: %v", err)
		} else {
			log.Printf("ERROR: unable to load config: %v", err)
		}
	} else {
	log.Print("Successfully loaded config.")
	}
}