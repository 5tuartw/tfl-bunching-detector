package stops

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/5tuartw/tfl-bunching-detector/internal/data"
	"github.com/5tuartw/tfl-bunching-detector/internal/helpers"
	"github.com/5tuartw/tfl-bunching-detector/internal/models"
)

// CSV headers:
// Stop_Code_LBSL,Bus_Stop_Code,Naptan_Atco,Stop_Name,Location_Easting,Location_Northing,Heading,Stop_Area,Virtual_Bus_Stop
const (
	naptanIdCol = 2
	stopNameCol = 3
	eastingCol  = 4
	northingCol = 5
	headingCol  = 6
	stopAreaCol = 7
)

func LoadBusStops() ([]models.BusStop, error) {

	busStopReader, err := data.NewBusStopReader()
	if err != nil {
		return nil, fmt.Errorf("could not open bus stop data: %v", err)
	}
	defer busStopReader.Close()

	checkDataFreshness(busStopReader.Info())

	csvReader := csv.NewReader(busStopReader)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv data: %w", err)
	}

	var allBusStops []models.BusStop

	stringConversionErrors := make(map[int][]string)

	for i, record := range records {
		if i == 0 {
			continue
		}
		easting, err := strconv.Atoi(record[eastingCol])
		if err != nil {
			easting = 0
			stringConversionErrors[i] = append(stringConversionErrors[i], "easting")
		}
		northing, err := strconv.Atoi(record[northingCol])
		if err != nil {
			northing = 0
			stringConversionErrors[i] = append(stringConversionErrors[i], "northing")
		}

		heading, err := strconv.Atoi(record[headingCol])
		if err != nil {
			heading = 0
		}

		allBusStops = append(allBusStops, models.BusStop{
			NaptanId:         record[naptanIdCol],
			StopName:         record[stopNameCol],
			LocationEasting:  easting,
			LocationNorthing: northing,
			Heading:          heading,
			StopArea:         record[stopAreaCol],
		})
	}

	if len(stringConversionErrors) > 0 {
		log.Printf("%d errors were found marshalling the csv data:", len(stringConversionErrors))
		for k, v := range stringConversionErrors {
			log.Printf("Line: %d - Errors: %v", k, v)
		}
	}

	return allBusStops, nil
}

func SearchStops(searchValue string, busStops []models.BusStop) []models.BusStop {
	var foundBusStops []models.BusStop

	for _, stop := range busStops {
		if strings.Contains(strings.ToLower(stop.StopName), strings.ToLower(searchValue)) {
			foundBusStops = append(foundBusStops, stop)
		}
	}

	return foundBusStops
}

func FindStopByID(naptainId string, busStops []models.BusStop) (models.BusStop, bool) {
	for _, stop := range busStops {
		if stop.NaptanId == naptainId {
			return stop, true
		}
	}
	return models.BusStop{}, false
}

func FindStopByIds(ids []string, allStops []models.BusStop) []models.BusStop {
	stopListWithDetails := []models.BusStop{}

	for _, stop := range ids {
		if stopInfo, found := FindStopByID(stop, allStops); found {
			stopListWithDetails = append(stopListWithDetails, stopInfo)
		} else {
			log.Printf("No stop found with Naptan ID '%s' whilst searching through route stops.", stop)
		}
	}
	return stopListWithDetails
}

func ChooseBusStop(busStops []models.BusStop) []models.BusStop {
	busStopCount := len(busStops)
	fmt.Printf("Found %d matching stops:\n", busStopCount)
	for i, stop := range busStops {
		fmt.Printf("%d: %s [%s] (%s)\n", i+1, stop.StopName, helpers.HeadingToDirection(stop.Heading), stop.NaptanId)
	}

	validEntry := false
	var userInput []int

	for !validEntry {
		userInput, validEntry = getStopInputHelper(busStopCount)
	}

	var chosenBusStops []models.BusStop

	for _, stop := range userInput {
		chosenBusStops = append(chosenBusStops, busStops[stop-1])
	}

	return chosenBusStops
}

func getStopInputHelper(maxChoice int) ([]int, bool) {
	var selectedStops []int
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter one or more stop numbers, separated by a space: ")

	if scanner.Scan() {
		userInput := scanner.Text()
		if userInput == "" {
			log.Println("Please enter at least one number.")
			return nil, false
		}
		selectedStopsStr := strings.Split(userInput, " ")
		for _, stopStr := range selectedStopsStr {
			stopInt, err := strconv.Atoi(stopStr)
			if err != nil {
				log.Printf("Error parsing choice: '%s'\n", stopStr)
				return nil, false
			}
			if stopInt > maxChoice || stopInt < 1 {
				log.Printf("Choice must be between 1-%d.\n", maxChoice)
				return nil, false
			}
			selectedStops = append(selectedStops, stopInt)
		}
	}

	return selectedStops, true
}
