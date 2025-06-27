package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/5tuartw/tfl-bunching-detector/internal/analysis"
	"github.com/5tuartw/tfl-bunching-detector/internal/config"
	"github.com/5tuartw/tfl-bunching-detector/internal/display"
	"github.com/5tuartw/tfl-bunching-detector/internal/helpers"
	"github.com/5tuartw/tfl-bunching-detector/internal/lines"
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
	}

	httpClient := tflclient.NewClient("https://api.tfl.gov.uk", cfg.TflKey)

	allBusStops, err := stops.LoadBusStops()
	if err != nil {
		log.Fatalf("Error: could not load bus stop data: %v", err)
	}

	// Analysing the given lineId at the current moment
	if *lineId != "" {
		fmt.Printf("Looking up routes and stops on bus line %s...\n\n", *lineId)
		lineInfo, err := lines.GetLineInfo(httpClient, *lineId)
		if err != nil {
			log.Fatalf("ERROR: could not get line information for %s: %v", *lineId, err)
		}

		selectedRoutes := lines.ChooseRoute(lineInfo)
		for _, route := range selectedRoutes {
			routeBunchingEvents := analysis.AnalyseRoute(*httpClient, *bunchingThreshold, lineInfo.Routes[route])
			routeName := lineInfo.Routes[route].Name
			display.PrintBunchingData(routeName, *bunchingThreshold, routeBunchingEvents)
		}
		os.Exit(1)
	}

	var chosenBusStops []models.BusStop

	// Conducting a search for stops if flag is raised and adding stopIds to list
	if *searchStop != "" {
		matchingBusStops := stops.SearchStops(*searchStop, allBusStops)
		if len(matchingBusStops) == 0 {
			log.Printf("No stops found containing '%s'.", *searchStop)
			os.Exit(0)
		}
		chosenBusStops = stops.ChooseBusStop(matchingBusStops)
	// otherwise adding the flagged stop id to the list of stops to list
	} else if *stopId != "" {
		if stop, found := stops.FindStopByID(*stopId, allBusStops); found {
			chosenBusStops = append(chosenBusStops, stop)
		} else {
			log.Fatalf("No stop found with Naptan ID '%s'.", *stopId)
		}
	}

	// looking up arrivalInfo for each selected stop and displaying to terminal
	for _, stop := range chosenBusStops {

		arrivalInfo, err := stops.GetStopArrivalInfo(httpClient, stop.NaptanId)
		if err != nil {
			log.Fatalf("ERROR: could not get arrival information for stop %s: %v", stop.NaptanId, err)
		}

		bunchingEvents := analysis.AnalyseArrivals(arrivalInfo, *bunchingThreshold)
		stopName := helpers.GetStopName(stop)
		display.PrintBunchingData(stopName, *bunchingThreshold, bunchingEvents)
	}
}
