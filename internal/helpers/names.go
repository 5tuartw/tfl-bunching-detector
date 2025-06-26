package helpers

import (
	"fmt"
	"strings"

	"github.com/5tuartw/tfl-bunching-detector/internal/models"
)

func GetStopName(stop models.BusStop) string {
	stopIdentifier := stop.NaptanId
	if stop.StopName != "" {
		stopIdentifier = fmt.Sprintf("%s [%s] (%s)", stop.StopName, HeadingToDirection(stop.Heading), stop.NaptanId)
	}
	return stopIdentifier
}

func CleanRouteName(name string) string {

	splitName := strings.Split(name, "&harr;")
	if len(splitName) != 2 {
		return name
	}
	rejoined := strings.Join(splitName, "to")
	removeSpaces := strings.Trim(strings.ReplaceAll(rejoined, "  ", " "), " ")
	return removeSpaces
}
