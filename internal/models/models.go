package models

import "time"

type Arrival struct {
	NaptanId      string `json:"naptanId"`
	StationName   string `json:"stationName"`
	LineId        string `json:"lineId"`
	LineName      string `json:"lineName"`
	VehicleId     string `json:"vehicleId"`
	TimeToStation int    `json:"timeToStation"`
}

type BunchingEvent struct {
	LineId      string    `json:"lineId"`
	NaptanId    string    `json:"naptanId"`
	StationName string    `json:"stationName"`
	EventTime   time.Time `json:"eventTime"`
	Headway     int       `json:"headway"`
	VehicleIds  []string  `json:"vehicleIds"`
}

type RouteBunches struct {
	RouteName string `json:"name"`
	
}

type BusStop struct {
	NaptanId         string `json:"naptanId"`
	StopName         string `json:"stopName"`
	LocationEasting  int    `json:"locationEasting"`
	LocationNorthing int    `json:"locationNorthing"`
	Heading          int    `json:"heading"`
	StopArea         string `json:"stopArea"`
}

type Route struct {
	Name        string    `json:"name"`
	StopIds     []string  `json:"naptanIds"`
	Stops       []BusStop `json:"busStops"`
	ServiceType string    `json:"serviceType"`
}

type LineInfo struct {
	LineId string  `json:"lineId"`
	Routes []Route `json:"orderedLineRoutes"`
}
