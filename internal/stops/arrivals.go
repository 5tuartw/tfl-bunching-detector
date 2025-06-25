package stops

import (
	"fmt"

	"github.com/5tuartw/tfl-bunching-detector/internal/models"
	"github.com/5tuartw/tfl-bunching-detector/internal/tflclient"
)

func GetStopArrivalInfo(httpClient *tflclient.Client, stopId string) ([]models.Arrival, error) {
	arrivalInfo, err := httpClient.RequestArrivalInfo(stopId)
	if err != nil {
		return nil, fmt.Errorf("could not request arrival info for %s via the TFL api: %v", stopId, err)
	}

	return arrivalInfo, nil
}
