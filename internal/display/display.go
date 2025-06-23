package display

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/5tuartw/tfl-bunching-detector/internal/helpers"
	"github.com/5tuartw/tfl-bunching-detector/internal/models"
)

func PrintBunchingData(stop models.BusStop, threshold int, bunchingEvents []models.BunchingEvent) {
	stopIdentifier := stop.NaptanId
	if stop.StopName != "" {
		stopIdentifier = fmt.Sprintf("%s [%s] (%s)", stop.StopName, helpers.HeadingToDirection(stop.Heading), stop.NaptanId)
	}
	if len(bunchingEvents) == 0 {
		log.Printf("No bus lines at stop %s are bunching right now.", stopIdentifier)
		return
	}

	fmt.Printf("\n%d bunching events found at stop %s (threshold: %d seconds):\n\n", len(bunchingEvents), stopIdentifier, threshold)

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Line\tVehicle 1\tVehicle 2\tHeadway (s)")
	fmt.Fprintln(tw, "====\t=========\t=========\t===========")
	for _, event := range bunchingEvents {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\n",
			event.LineId,
			event.VehicleIds[0],
			event.VehicleIds[1],
			event.Headway,
		)
	}

	tw.Flush()

}
