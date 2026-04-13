package models

import (
	"encoding/json"
	"fmt"
)

type TelegramEventData struct {
	EventType       string `json:"event_type"`
	ChatID          int64  `json:"chat_id"`
	ChannelTitle    string `json:"channel_title"`
	ChannelMainLang string `json:"channel_main_lang"`
	MsgID           int    `json:"msg_id"`
}

type RawTelegramEvent struct {
	rawOsintEvent
	Data TelegramEventData `json:"data"`
}

func (e *RawTelegramEvent) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, e); err != nil {
		return fmt.Errorf("RawTelegramEvent unmarshal error: %w", err)
	}
	return nil
}
