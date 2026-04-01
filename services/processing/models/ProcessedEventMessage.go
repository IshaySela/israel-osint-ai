package models

type ProcessedEventMessage struct {
	EsId      string
	Summary   string             `json:"summary"`
	Locations map[string]Geocode `json:"locations"`
	Timestamp string
}
