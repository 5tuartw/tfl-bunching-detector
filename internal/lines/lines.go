package lines

import (
	"fmt"

	"github.com/5tuartw/tfl-bunching-detector/internal/helpers"
	"github.com/5tuartw/tfl-bunching-detector/internal/models"
	"github.com/5tuartw/tfl-bunching-detector/internal/tflclient"
)

func GetLineInfo(client *tflclient.Client, lineId string) (models.LineInfo, error) {
	lineInfo, err := client.RequestLineInfo(lineId)
	if err != nil {
		return models.LineInfo{}, err
	}
	for i, route := range lineInfo.Routes {
		lineInfo.Routes[i].Name = helpers.CleanRouteName(route.Name)
	}
	return lineInfo, nil
}

func ChooseRoute(line models.LineInfo) []int {
	routeCount := len(line.Routes)
	fmt.Printf("Found %d routes on line %s:\n", routeCount, line.LineId)
	for i, route := range line.Routes {
		fmt.Printf("%d: %s\n", i+1, route.Name)
	}

	validEntry := false
	var userInput []int

	for !validEntry {
		userInput, validEntry = helpers.GetMenuChoiceMultiple(routeCount)
	}

	var chosenRoutes []int

	for _, route := range userInput {
		chosenRoutes = append(chosenRoutes, route-1)
	}

	return chosenRoutes
}
