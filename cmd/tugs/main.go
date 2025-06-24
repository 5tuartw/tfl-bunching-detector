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

	stopId := flag.String("stop-id", "", "NaptanId for a specific stop")
	bunchingThreshold := flag.Int("threshold", 90, "threshold for bunched buses in seconds")
	searchStop := flag.String("search", "", "search for a stop by name")
	lineId := flag.String("line", "", "Line ID or Line Number")

	flag.Parse()

	if *searchStop == "" && *stopId == "" && *lineId == "" {
		log.Fatal("Please use an appropriate flag to select stops(s) or line.")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("ERROR: unable to load config: %v", err)
	} else {
		log.Printf("Successfully loaded config. Using API key starting with: %s...", cfg.TflKey[:4])
	}

	httpClient := tflclient.NewClient("https://api.tfl.gov.uk", cfg.TflKey)

	allBusStops, err := stops.LoadBusStops()
	if err != nil {
		log.Fatalf("Error: could not load bus stop data: %v", err)
	}

	if *lineId != "" {
		lineInfo, err := httpClient.GetLineInfo(*lineId)
		if err != nil {
			log.Fatalf("ERROR: could not get line information for %s: %v", *lineId, err)
		}

		for i, route := range lineInfo.Routes {
			stopListComplete := []models.BusStop{}

			for _, stop := range route.StopIds {
				if stopInfo, found := stops.FindStopByID(stop, allBusStops); found {
					stopListComplete = append(stopListComplete, stopInfo)
				} else {
					log.Printf("No stop found with Naptan ID '%s'.", stop)
				}
			}
			lineInfo.Routes[i].Stops = stopListComplete
			display.PrintRoute(lineInfo.Routes[i])
		}
		os.Exit(1)
	}

	var chosenBusStops []models.BusStop

	if *searchStop != "" {
		matchingBusStops := stops.SearchStops(*searchStop, allBusStops)
		if len(matchingBusStops) == 0 {
			log.Printf("No stops found containing '%s'.", *searchStop)
			os.Exit(0)
		}
		chosenBusStops = stops.ChooseBusStop(matchingBusStops)

	} else if *stopId != "" {
		if stop, found := stops.FindStopByID(*stopId, allBusStops); found {
			chosenBusStops = append(chosenBusStops, stop)
		} else {
			log.Fatalf("No stop found with Naptan ID '%s'.", *stopId)
		}
	}

	for _, stop := range chosenBusStops {

		arrivalInfo, err := httpClient.GetArrivalInfo(stop.NaptanId)
		if err != nil {
			log.Fatalf("ERROR: could not get arrival information for stop %s: %v", stop.NaptanId, err)
		}

		bunchingEvents := analysis.AnalyseArrivals(arrivalInfo, *bunchingThreshold)
		display.PrintBunchingData(stop, *bunchingThreshold, bunchingEvents)
	}
}
