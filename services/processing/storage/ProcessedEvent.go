package storage

import (
	models "processing/models"
)

type ProcessedEvent[T any] struct {
	RawMessage     string            `json:"raw_message"`
	Summary        string            `json:"summary"`
	Locations      []models.Location `json:"locations"`
	TimestampEpoch int64             `json:"timestamp_epoch"`
	Source         string            `json:"source"`
	Data           T                 `json:"data"`
}

type ProcessedTelegramEvent = ProcessedEvent[models.TelegramEventData]
