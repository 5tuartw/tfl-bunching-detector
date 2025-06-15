package models

type Arrival struct {
	NaptanId      string `json:"naptanId"`
	StationName   string `json:"stationName"`
	LineId        string `json:"lineId"`
	LineName      string `json:"lineName"`
	VehicleId     string `json:"vehicleId"`
	TimeToStation int    `json:"timeToStation"`
}
