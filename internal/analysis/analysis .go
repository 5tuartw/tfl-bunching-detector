package analysis

import (
	"sort"
	"time"

	"github.com/5tuartw/tfl-bunching-detector/internal/models"
)

func AnalyseArrivals(arrivals []models.Arrival, bunchingThreshold int) []models.BunchingEvent {
	groupedArrivals := groupByLine(arrivals)
	var bunchingEvents []models.BunchingEvent

	for line := range groupedArrivals {
		sort.Slice(groupedArrivals[line], func(i, j int) bool {
			return groupedArrivals[line][i].TimeToStation < groupedArrivals[line][j].TimeToStation
		})

		for idx := 0; idx < len(groupedArrivals[line])-1; idx++ {
			headway := groupedArrivals[line][idx+1].TimeToStation - groupedArrivals[line][idx].TimeToStation
			if headway < bunchingThreshold {
				bunchingEvents = append(bunchingEvents, models.BunchingEvent{
					LineId:      line,
					NaptanId:    groupedArrivals[line][idx].NaptanId,
					StationName: groupedArrivals[line][idx].StationName,
					EventTime:   time.Now(),
					Headway:     headway,
					VehicleIds:  []string{groupedArrivals[line][idx].VehicleId, groupedArrivals[line][idx+1].VehicleId},
				})
			}
		}
	}

	return bunchingEvents
}

func groupByLine(a []models.Arrival) map[string][]models.Arrival {

	groupedArrivals := make(map[string][]models.Arrival)

	for _, arrival := range a {
		// append will automatically create a new entry if arrival.LineId is not already in the map
		// due to maps returning nil when the key is not found, and append(nil, item) will create a new slice
		groupedArrivals[arrival.LineId] = append(groupedArrivals[arrival.LineId], arrival)
	}

	return groupedArrivals
}
