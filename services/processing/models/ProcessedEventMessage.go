package models

type ProcessedEventMessage struct {
	DbId      string     `json:"dbId"`
	Summary   string     `json:"summary"`
	Locations []Location `json:"locations"`
	Timestamp string     `json:"timestamp"`
}
