package modelPostgresql

import "time"

type PowerMeter struct {
	Time      time.Time
	DeviceID  string
	Metric    string
	Value     float64
	EventHash string
}
