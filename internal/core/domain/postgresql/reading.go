package modelPostgresql

import "time"

type PowerMeter struct {
	Time      time.Time `json:"time"`
	DeviceID  string    `json:"device_id"`
	Metric    string    `json:"metric"`
	Value     float64   `json:"value"`
	EventHash string    `json:"event_hash"`
}
