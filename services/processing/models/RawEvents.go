package models

import (
	"encoding/json"
	"fmt"
)

type RawTelegramEvent struct {
	rawOsintEvent
	Text            string `json:"text"`
	EventType       string `json:"event_type"`
	ChatID          int64  `json:"chat_id"`
	ChannelTitle    string `json:"channel_title"`
	ChannelMainLang string `json:"channel_main_lang"`
	MessageID       int    `json:"message_id"`
}

func (e *RawTelegramEvent) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, e); err != nil {
		return fmt.Errorf("RawTelegramEvent unmarshal error: %w", err)
	}
	return nil
}
