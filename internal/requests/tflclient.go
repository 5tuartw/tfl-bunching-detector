package tflclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/5tuartw/tfl-bunching-detector/internal/models"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetArrivalInfo(stopId string) ([]models.Arrival, error) {
	arrivalsRequestUrl := c.BaseURL + "/StopPoint/" + stopId + "/arrivals"
	var arrivalInfo []models.Arrival
	var requestBody io.Reader

	arrivalInfoRequest, err := http.NewRequest("GET", arrivalsRequestUrl, requestBody)
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %v", err)
	}

	response, err := c.HTTPClient.Do(arrivalInfoRequest)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	if response == nil {
		return nil, fmt.Errorf("response from arrivals request was nil")
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code received %d: %s", response.StatusCode, response.Status)
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response data: %v", err)
	}

	err = json.Unmarshal(data, &arrivalInfo)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json arrival data: %v", err)
	}

	return arrivalInfo, nil
}
