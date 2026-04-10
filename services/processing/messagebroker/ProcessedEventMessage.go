package messagebroker

import (
	models "github.com/IshaySela/israel-osint-ai/services/processing/models"
	storage "github.com/IshaySela/israel-osint-ai/services/processing/storage"
)

type EventsMessageFields struct {
	DbId   string `json:"dbId"`
	Source string `json:"source"`
}

type ProcessedEventMessage interface{}

type ProcessedTelegramEvMessage struct {
	EventsMessageFields
	Summary         string            `json:"summary"`
	Locations       []models.Location `json:"locations"`
	TimestampEpoch  int64             `json:"timestamp_epoch"`
	ChannelTitle    string            `json:"channel_title"`
	ChannelMainLang string            `json:"channel_main_lang"`
}

// CreateMessageFromEvent recives an event of some types and creates a message for the message
// broker. returns nil if ievent is of invalid type.
func CreateMessageFromEvent(ievent storage.IProcessedEvent, dbId string) ProcessedEventMessage {
	var result ProcessedEventMessage

	switch event := ievent.(type) {
	case storage.ProcessedTelegramEvent:
		result = ProcessedTelegramEvMessage{
			EventsMessageFields: EventsMessageFields{
				DbId:   dbId,
				Source: event.Source,
			},
			Summary:         event.Summary,
			Locations:       event.Locations,
			TimestampEpoch:  event.TimestampEpoch,
			ChannelTitle:    event.ChannelTitle,
			ChannelMainLang: event.ChannelMainLang,
		}
	default:
		result = nil
	}

	return result
}
