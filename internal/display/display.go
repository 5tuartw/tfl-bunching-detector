package display

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/5tuartw/tfl-bunching-detector/internal/helpers"
	"github.com/5tuartw/tfl-bunching-detector/internal/models"
)

func PrintBunchingData(name string, threshold int, bunchingEvents []models.BunchingEvent) {

	if len(bunchingEvents) == 0 {
		log.Printf("No buses are bunching right now (%s).", name)
		return
	}

	fmt.Printf("\n%d bunching events found for %s (threshold: %d seconds):\n\n", len(bunchingEvents), name, threshold)

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Line\tStop\tVehicle 1\tVehicle 2\tHeadway (s)")
	fmt.Fprintln(tw, "====\t====\t=========\t=========\t===========")
	for _, event := range bunchingEvents {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\n",
			event.LineId,
			event.StationName,
			event.VehicleIds[0],
			event.VehicleIds[1],
			event.Headway,
		)
	}

	tw.Flush()

}

func PrintRoute(route models.Route) {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Printf("Route: %s\n\n", route.Name)
	fmt.Fprintln(tw, "Stop\tHeading")
	fmt.Fprintln(tw, "====\t=======")
	for _, stop := range route.Stops {
		fmt.Fprintf(tw, "%s\t%s\n",
			stop.StopName,
			helpers.HeadingToDirection(stop.Heading),
		)
	}
	tw.Flush()
	fmt.Println("")
}
