package main

import (
	"fmt"
	"log"
	"os"

	"github.com/5tuartw/tfl-bunching-detector/internal/analysis"
	"github.com/5tuartw/tfl-bunching-detector/internal/config"
	"github.com/5tuartw/tfl-bunching-detector/internal/display"
	"github.com/5tuartw/tfl-bunching-detector/internal/helpers"
	"github.com/5tuartw/tfl-bunching-detector/internal/models"
	"github.com/5tuartw/tfl-bunching-detector/internal/stops"
	"github.com/5tuartw/tfl-bunching-detector/internal/tflclient"
	"github.com/5tuartw/tfl-bunching-detector/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("ERROR: unable to load config: %v", err)
	}
	tflClient := tflclient.NewClient("https://api.tfl.gov.uk", cfg.TflKey)

	allStops, err := stops.LoadBusStops()
	if err != nil {
		log.Fatalf("Failed to load bus stops: %v", err)
	}

	tuiModel := tui.NewModel(allStops, tflClient)

	p := tea.NewProgram(tuiModel)
	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	m, ok := finalModel.(*tui.Model)
	if !ok {
		os.Exit(1)
	}

	if len(m.SelectedStops) > 0 {
		fmt.Println("\n--- Analysing Selected Stop ---")
		var chosenStops []models.BusStop
		for index := range m.SelectedStops {
			chosenStops = append(chosenStops, m.SearchResults[index])
		}

		for _, stop := range chosenStops {
			fmt.Printf("\nFetching data for stop: %s\n", stop.StopName)
			arrivalInfo, err := stops.GetStopArrivalInfo(tflClient, stop.NaptanId)
			if err != nil {
				log.Fatalf("EROR: could not get arrival information for stop %s (%s): %v", stop.StopName, stop.NaptanId, err)
			}
			bunchingEvents := analysis.AnalyseArrivals(arrivalInfo, "", m.Threshold)
			stopName := helpers.GetStopName(stop)
			display.PrintBunchingData(stopName, m.Threshold, bunchingEvents)
		}
	}
}
