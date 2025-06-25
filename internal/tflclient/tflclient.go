package tflclient

import (
	"encoding/json"
	"fmt"
	"io"
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

// helper function to send an http request and marshall into the target interface model
func (c *Client) doRequest(req *http.Request, target interface{}) error {

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	if response == nil {
		return fmt.Errorf("response from request was nil")
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("non-200 status code received %d: %s", response.StatusCode, response.Status)
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("could not read response data: %v", err)
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("could not unmarshal json data: %v", err)
	}
	return nil
}

func (c *Client) RequestArrivalInfo(stopId string) ([]models.Arrival, error) {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("base url (%s) could not be parsed in GetArrivalInfo: %v", c.BaseURL, err)
	}
	baseURL = baseURL.JoinPath("StopPoint", stopId, "arrivals")
	q := baseURL.Query()
	q.Add("app_key", c.APIKey)
	baseURL.RawQuery = q.Encode()
	finalURL := baseURL.String()

	var arrivalInfo []models.Arrival
	var requestBody io.Reader

	arrivalInfoRequest, err := http.NewRequest("GET", finalURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("could not create http request: %v", err)
	}
	arrivalInfoRequest.Header.Set("User-Agent", "TFL-Bus-Bunching-Detector/1.0 stuart@stuartw.dev")

	err = c.doRequest(arrivalInfoRequest, &arrivalInfo)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %v", err)
	}

	return arrivalInfo, nil
}

func (c *Client) RequestLineInfo(lineId string) (models.LineInfo, error) {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return models.LineInfo{}, fmt.Errorf("base url (%s) could not be parsed in GetLineInfo: %v", c.BaseURL, err)
	}
	baseURL = baseURL.JoinPath("Line", lineId, "Route/Sequence/All")
	q := baseURL.Query()
	q.Add("app_key", c.APIKey)
	baseURL.RawQuery = q.Encode()
	finalURL := baseURL.String()

	var lineInfo models.LineInfo
	var requestBody io.Reader

	lineInfoRequest, err := http.NewRequest("GET", finalURL, requestBody)
	if err != nil {
		return models.LineInfo{}, fmt.Errorf("could not create http request: %v", err)
	}
	lineInfoRequest.Header.Set("User-Agent", "TFL-Bus-Bunching-Detector/1.0 stuart@stuartw.dev")

	err = c.doRequest(lineInfoRequest, &lineInfo)
	if err != nil {
		return models.LineInfo{}, fmt.Errorf("http request failed: %v", err)
	}

	return lineInfo, nil
}
