package tflclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("base url (%s) could not be parsed in GetArrivalInfo: %v", c.BaseURL, err)
	}
	baseURL = baseURL.JoinPath("StopPoint", stopId, "arrivals")
	q := baseURL.Query()
	q.Add("app_key", c.APIKey)
	baseURL.RawQuery = q.Encode()
	finalURL := baseURL.String()
	log.Print(finalURL)

	var arrivalInfo []models.Arrival
	var requestBody io.Reader

	arrivalInfoRequest, err := http.NewRequest("GET", finalURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %v", err)
	}
	arrivalInfoRequest.Header.Set("User-Agent", "TFL-Bus-Bunching-Detector/1.0 stuart@stuartw.dev")
	
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
