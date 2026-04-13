package messagebroker

import (
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	storage "github.com/IshaySela/israel-osint-ai/services/processing/storage"
)

type EventsMessageFields struct {
}

type ProcessedEventMessage[T any] struct {
	DbId           string            `json:"dbId"`
	Source         string            `json:"source"`
	Summary        string            `json:"summary"`
	Locations      []models.Location `json:"locations"`
	TimestampEpoch int64             `json:"timestamp_epoch"`
	Data           T                 `json:"data"`
}

type ProcessedTelegramEvMessage = ProcessedEventMessage[models.TelegramEventData]

// CreateMessageFromEvent recives an event of some types and creates a message for the message
// broker. returns nil if ievent is of invalid type.
func CreateMessageFromEvent(event storage.ProcessedEvent[any], dbId string) ProcessedEventMessage[any] {
	return ProcessedEventMessage[any]{
		DbId:           dbId,
		Source:         event.Source,
		Summary:        event.Summary,
		Locations:      event.Locations,
		TimestampEpoch: event.TimestampEpoch,
		Data:           event.Data,
	}
}
