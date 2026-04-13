package storage

import (
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
)

type IProcessedEvent interface {
	// the empty interface is used in place where a func expects either type
	// of processed event
}

type ProcessedEvent struct {
	RawMessage     string            `json:"raw_message"`
	Summary        string            `json:"summary"`
	Locations      []models.Location `json:"locations"`
	TimestampEpoch int64             `json:"timestamp_epoch"`
	Source         string            `json:"source"`
}

type ProcessedTelegramEvent struct {
	ProcessedEvent
	Data interface{} `json:"data"`
}
