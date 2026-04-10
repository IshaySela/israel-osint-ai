package models

type ProcessedEventMessage struct {
	DbId            string     `json:"dbId"`
	Summary         string     `json:"summary"`
	Locations       []Location `json:"locations"`
	TimestampEpoch  int64      `json:"timestamp_epoch"`
	ChannelTitle    string     `json:"channel_title"`
	ChannelMainLang string     `json:"channel_main_lang"`
}
