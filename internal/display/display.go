package display

import (
	"fmt"
	"log"

	"github.com/5tuartw/tfl-bunching-detector/internal/models"
)

func PrintBunchingData(stop models.BusStop, threshold int, bunchingEvents []models.BunchingEvent) {
	stopIdentifier := stop.NaptanId
	if stop.StopName != "" {
		stopIdentifier = fmt.Sprintf("%s (%s)", stop.StopName, stop.NaptanId)
	}
	if len(bunchingEvents) == 0 {
		log.Printf("No bus lines at stop %s are bunching right now.", stopIdentifier)
		return
	}
	fmt.Printf("%d bunching events found at stop %s (threshold: %d seconds):\n", len(bunchingEvents), stopIdentifier, threshold)
	for i, event := range bunchingEvents {
		fmt.Printf("%d. Line %s: Vehicles: %s and %s are %d seconds apart.\n", i+1, event.LineId, event.VehicleIds[0], event.VehicleIds[1], event.Headway)
	}
}
