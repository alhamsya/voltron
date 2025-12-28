package util

import (
	"fmt"
	"strings"
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

func ParseDeviceTime(s string, loc *time.Location) (time.Time, error) {
	const (
		layoutA = "02/01/2006 15:04:05" // contoh lama: 19/12/2025 15:27:53
		layoutB = "02 15:04:05/01/2006" // contoh kamu: 28 10:54:15/12/2025
	)

	s = strings.TrimSpace(s)

	if t, err := time.ParseInLocation(layoutB, s, loc); err == nil {
		return t, nil
	}

	return time.ParseInLocation(layoutA, s, loc)
}
