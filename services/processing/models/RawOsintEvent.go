package models

import (
	"encoding/json"
)

type RawOsintEvent struct {
	AppEvId    string          `json:"app_ev_id"`
	RawMessage string          `json:"raw_message"`
	Timestamp  string          `json:"timestamp"`
	Source     string          `json:"source"`
	Data       json.RawMessage `json:"data"`
}

func ParseRawOsintEvent(b []byte) (RawOsintEvent, error) {
	var parsed RawOsintEvent
	err := json.Unmarshal(b, &parsed)

	if err != nil {
		return RawOsintEvent{}, err
	}

	return parsed, nil
}
