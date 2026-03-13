package models

import (
	"encoding/json"
	"fmt"
)

type RawOsintEvent struct {
	Text      string `json:"text"`
	EventType string `json:"event_type"`
	ChatID    int64  `json:"chat_id"`
	MessageID int    `json:"message_id"`
	Date      string `json:"date"`
}

func (e *RawOsintEvent) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, e); err != nil {
		return fmt.Errorf("RawOsintEvent unmarshal error: %w", err)
	}
	return nil
}
