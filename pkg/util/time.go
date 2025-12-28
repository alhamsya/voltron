package util

import (
	"fmt"
	"time"
)

func FormatMinutes(m int) string {
	h := m / 60
	minResult := m % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, minResult)
	}
	return fmt.Sprintf("%dm", minResult)
}

func WithinTimeWindow(t time.Time, from time.Time, to time.Time) bool {
	fm := from.Hour()*60 + from.Minute()
	tm := to.Hour()*60 + to.Minute()
	cur := t.Hour()*60 + t.Minute()

	// normal window (08:00–17:00)
	if fm <= tm {
		return cur >= fm && cur <= tm
	}
	// wrap window (22:00–02:00)
	return cur >= fm || cur <= tm
}
