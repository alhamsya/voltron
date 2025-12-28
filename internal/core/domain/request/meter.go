package modelRequest

import "time"

type ReqHandlerMeterReading struct {
	PowerMeter []HandlerPowerMater `json:"PM"`
}

type HandlerPowerMater struct {
	Date string `json:"date"`
	Data string `json:"data"`
	Name string `json:"name"`
}

type PowerMater struct {
	Date time.Time
	Data float64
	Name string
}
