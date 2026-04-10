package models

import "time"

// ParseToEpoch parses a datetime string (Python str(datetime) format) to Unix epoch seconds.
// Falls back to time.Now().Unix() if parsing fails.
func ParseToEpoch(s string) int64 {
	if t, err := time.Parse("2006-01-02 15:04:05-07:00", s); err == nil {
		return t.Unix()
	}
	return time.Now().Unix()
}
