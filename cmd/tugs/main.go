package main

import (
	"flag"
	"log"
	"os"

	"github.com/5tuartw/tfl-bunching-detector/internal/analysis"
	"github.com/5tuartw/tfl-bunching-detector/internal/config"
	"github.com/5tuartw/tfl-bunching-detector/internal/display"
	"github.com/5tuartw/tfl-bunching-detector/internal/models"
	"github.com/5tuartw/tfl-bunching-detector/internal/stops"
	"github.com/5tuartw/tfl-bunching-detector/internal/tflclient"
)

func main() {

	stopId := flag.String("stop-id", "490000234H", "NaptanId for a specific stop")
	bunchingThreshold := flag.Int("threshold", 90, "threshold for bunched buses in seconds")
	searchStop := flag.String("search", "", "search for a stop by name")
	flag.Parse()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("ERROR: unable to load config: %v", err)
	} else {
		log.Printf("Successfully loaded config. Using API key starting with: %s...", cfg.TflKey[:4])
	}

	httpClient := tflclient.NewClient("https://api.tfl.gov.uk", cfg.TflKey)

	var chosenBusStops []models.BusStop

	if *searchStop != "" {
		allBusStops, err := stops.LoadBusStops()
		if err != nil {
			log.Fatalf("Error: could not load bus stop data: %v", err)
		}
		matchingBusStops := stops.SearchStops(*searchStop, allBusStops)
		if len(matchingBusStops) == 0 {
			log.Printf("No stops found containing '%s'.", *searchStop)
			os.Exit(0)
		}
		chosenBusStops = stops.ChooseBusStop(matchingBusStops)

	} else if *stopId != "" {
			chosenBusStops = append(chosenBusStops, models.BusStop{
				NaptanId: *stopId,
			})
	} else {
		log.Fatal("Please use the -search or -stop-id flags to select stop(s).")
	}

	for _, stop := range chosenBusStops {

		arrivalInfo, err := httpClient.GetArrivalInfo(stop.NaptanId)
		if err != nil {
			log.Fatalf("ERROR: could not get arrival information for stop %s: %v", stop.NaptanId, err)
		}

		bunchingEvents := analysis.AnalyseArrivals(arrivalInfo, *bunchingThreshold)

		display.PrintBunchingData(stop.NaptanId, *bunchingThreshold, bunchingEvents)
	}
}
