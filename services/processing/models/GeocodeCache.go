package models

type GeocodeCache struct {
	LocationText string  `json:"location_text"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	Timestamp    string  `json:"timestamp"`
}
