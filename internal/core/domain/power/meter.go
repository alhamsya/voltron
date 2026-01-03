package modelPower

type BillingSummary struct {
	DeviceID   string  `json:"device_id"`
	From       string  `json:"from"`
	To         string  `json:"to"`
	TotalKwh   float64 `json:"total_kwh"`
	RatePerKwh float64 `json:"rate_per_kwh"`
	Subtotal   float64 `json:"subtotal"`
	Tax        float64 `json:"tax"`
	Total      float64 `json:"total"`
}

type DailyLine struct {
	Day      string  `json:"day"` // YYYY-MM-DD
	UsageKwh float64 `json:"usage_kwh"`
}
