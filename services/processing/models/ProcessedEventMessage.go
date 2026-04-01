package models

type ProcessedEventMessage struct {
	DbId      string             `json:"dbId"`
	Summary   string             `json:"summary"`
	Locations map[string]Geocode `json:"locations"`
	Timestamp string             `json:"timestamp"`
}
