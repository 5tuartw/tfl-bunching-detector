package main

import (
	"log"
	"os"

	"github.com/5tuartw/tfl-bunching-detector/internal/analysis"
	"github.com/5tuartw/tfl-bunching-detector/internal/config"
	"github.com/5tuartw/tfl-bunching-detector/internal/tflclient"
)

const bunchingThreshold = 600 // seconds

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Printf("ERROR: unable to load config: %v", err)
		os.Exit(1)
	} else {
		log.Printf("Successfully loaded config. Using API key starting with: %s...", cfg.TflKey[:4])
	}

	stopId := "490000234H"
	httpClient := tflclient.NewClient("https://api.tfl.gov.uk", cfg.TflKey)

	arrivalInfo, err := httpClient.GetArrivalInfo(stopId)
	if err != nil {
		log.Printf("ERROR: could not get arrival information for stop %s: %v", stopId, err)
		os.Exit(1)
	}

	bunchingEvents := analysis.AnalyseArrivals(arrivalInfo, bunchingThreshold)

	log.Print(arrivalInfo)
	log.Print(bunchingEvents)
}
