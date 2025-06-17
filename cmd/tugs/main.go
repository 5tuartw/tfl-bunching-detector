package main

import (
	"flag"
	"log"

	"github.com/5tuartw/tfl-bunching-detector/internal/analysis"
	"github.com/5tuartw/tfl-bunching-detector/internal/config"
	"github.com/5tuartw/tfl-bunching-detector/internal/display"
	"github.com/5tuartw/tfl-bunching-detector/internal/tflclient"
)

func main() {

	stopId := flag.String("stop-id", "490000234H", "NaptanId for a specific stop")
	bunchingThreshold := flag.Int("threshold", 90, "threshold for bunched buses in seconds")
	flag.Parse()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("ERROR: unable to load config: %v", err)
	} else {
		log.Printf("Successfully loaded config. Using API key starting with: %s...", cfg.TflKey[:4])
	}

	httpClient := tflclient.NewClient("https://api.tfl.gov.uk", cfg.TflKey)

	arrivalInfo, err := httpClient.GetArrivalInfo(*stopId)
	if err != nil {
		log.Fatalf("ERROR: could not get arrival information for stop %s: %v", *stopId, err)
	}

	bunchingEvents := analysis.AnalyseArrivals(arrivalInfo, *bunchingThreshold)

	display.PrintBunchingData(*stopId, *bunchingThreshold, bunchingEvents)
}
