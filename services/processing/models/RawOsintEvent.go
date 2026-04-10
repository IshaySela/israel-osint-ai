package models

import (
	"encoding/json"
	"fmt"
)

type RawOsintEvent interface{}

type rawOsintEvent struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
}

func ParseRawOsintEvent(b []byte) (RawOsintEvent, error) {
	var raw rawOsintEvent
	var parsed RawOsintEvent
	err := json.Unmarshal(b, &raw)

	if err != nil {
		return nil, err
	}

	switch raw.Source {
	case "telegram":
		var tg RawTelegramEvent
		err = json.Unmarshal(b, &tg)

		if err != nil {
			return nil, err
		}

		parsed = tg
	default:
		return nil, fmt.Errorf("Unknown source %s", raw.Source)
	}

	return parsed, nil
}
