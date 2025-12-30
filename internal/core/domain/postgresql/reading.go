package modelPostgresql

import "time"

type PowerMeter struct {
	Time      time.Time `json:"time"`
	DeviceID  string    `json:"device_id"`
	Metric    string    `json:"metric"`
	Value     float64   `json:"value"`
	EventHash string    `json:"event_hash"`
}

type TimeSeries struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

type Latest struct {
	Metric string    `json:"metric"`
	Time   time.Time `json:"time"`
	Value  float64   `json:"value"`
}

type DailyUsage struct {
	Day      string  `json:"day"` // "2025-12-19"
	UsageKwh float64 `json:"usage_kwh"`
}
