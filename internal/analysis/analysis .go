package analysis

import (
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/5tuartw/tfl-bunching-detector/internal/models"
	"github.com/5tuartw/tfl-bunching-detector/internal/stops"
	"github.com/5tuartw/tfl-bunching-detector/internal/tflclient"
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

func groupByLine(arrivals []models.Arrival) map[string][]models.Arrival {

	groupedArrivals := make(map[string][]models.Arrival)

	for _, arrival := range arrivals {
		// append will automatically create a new entry if arrival.LineId is not already in the map
		// due to maps returning nil when the key is not found, and append(nil, item) will create a new slice
		groupedArrivals[arrival.LineId] = append(groupedArrivals[arrival.LineId], arrival)
	}

	return groupedArrivals
}

func AnalyseRoute(client tflclient.Client, threshold int, route models.Route) []models.BunchingEvent {
	bunchesOnRoute := []models.BunchingEvent{}

	var wg sync.WaitGroup
	jobs := make(chan string, len(route.StopIds))
	results := make(chan []models.BunchingEvent)

	// setup the workers
	const numWorkers = 5
	for w := 1; w <= numWorkers; w++ {
		go routeWorker(&wg, &client, threshold, jobs, results)
	}

	//assign the workers
	for _, stop := range route.StopIds {
		wg.Add(1)
		jobs <- stop
	}
	close(jobs)

	//prepare to close results when done
	go func() {
		wg.Wait()
		close(results)
	}()

	//add results to list
	for res := range results {
		bunchesOnRoute = append(bunchesOnRoute, res...)
	}

	bunchesWithoutRepeats := removeRepeats(bunchesOnRoute)
	return bunchesWithoutRepeats
}

func routeWorker(wg *sync.WaitGroup, client *tflclient.Client, threshold int, jobs <-chan string, results chan<- []models.BunchingEvent) {
	for stopId := range jobs {
		arrivals, err := stops.GetStopArrivalInfo(client, stopId)
		if err != nil {
			log.Printf("Unable to get arrival info for stop '%s': %v", stopId, err)
		}
		bunches := AnalyseArrivals(arrivals, threshold)
		results <- bunches
		wg.Done()
	}
}

func removeRepeats(bunches []models.BunchingEvent) []models.BunchingEvent {
	cleanedBunches := []models.BunchingEvent{}
	seen := make(map[string]bool)

	for _, bunch := range bunches {
		sortedVehicles := make([]string, len(bunch.VehicleIds))
		_ = copy(sortedVehicles, bunch.VehicleIds)
		sort.Strings(sortedVehicles)
		key := strings.Join(sortedVehicles, "")
		if !seen[key] {
			cleanedBunches = append(cleanedBunches, bunch)
			seen[key] = true
		}
	}

	return cleanedBunches
}
