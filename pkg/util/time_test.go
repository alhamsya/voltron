package util

import (
	"testing"
	"time"
)

func TestFormatMinutes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   int
		want string
	}{
		{"zero", 0, "0m"},
		{"minutes only", 5, "5m"},
		{"exact 1 hour", 60, "1h 0m"},
		{"1 hour 1 minute", 61, "1h 1m"},
		{"2 hours 30 minutes", 150, "2h 30m"},
		{"59 minutes", 59, "59m"},
		{"large", 24*60 + 15, "24h 15m"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := FormatMinutes(tt.in); got != tt.want {
				t.Fatalf("FormatMinutes(%d) = %q; want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestWithinTimeWindow(t *testing.T) {
	t.Parallel()

	// Helper bikin time dengan jam:menit (date irrelevant karena fungsi cuma pakai hour/minute)
	at := func(h, m int) time.Time {
		return time.Date(2025, 1, 1, h, m, 0, 0, time.UTC)
	}

	tests := []struct {
		name string
		t    time.Time
		from time.Time
		to   time.Time
		want bool
	}{
		// Normal window (08:00–17:00)
		{"normal: inside mid", at(10, 0), at(8, 0), at(17, 0), true},
		{"normal: on start boundary", at(8, 0), at(8, 0), at(17, 0), true},
		{"normal: on end boundary", at(17, 0), at(8, 0), at(17, 0), true},
		{"normal: before start", at(7, 59), at(8, 0), at(17, 0), false},
		{"normal: after end", at(17, 1), at(8, 0), at(17, 0), false},

		// Wrap window (22:00–02:00)
		{"wrap: inside before midnight", at(23, 0), at(22, 0), at(2, 0), true},
		{"wrap: inside after midnight", at(1, 0), at(22, 0), at(2, 0), true},
		{"wrap: on start boundary", at(22, 0), at(22, 0), at(2, 0), true},
		{"wrap: on end boundary", at(2, 0), at(22, 0), at(2, 0), true},
		{"wrap: outside (middle day)", at(12, 0), at(22, 0), at(2, 0), false},
		{"wrap: just outside after end", at(2, 1), at(22, 0), at(2, 0), false},
		{"wrap: just outside before start", at(21, 59), at(22, 0), at(2, 0), false},

		// Edge case: from == to (window 1 titik waktu) -> dengan logic kamu masuk ke “normal window”
		{"same from/to: exact match", at(9, 0), at(9, 0), at(9, 0), true},
		{"same from/to: different time", at(9, 1), at(9, 0), at(9, 0), false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := WithinTimeWindow(tt.t, tt.from, tt.to); got != tt.want {
				t.Fatalf(
					"WithinTimeWindow(t=%s, from=%s, to=%s) = %v; want %v",
					tt.t.Format("15:04"), tt.from.Format("15:04"), tt.to.Format("15:04"),
					got, tt.want,
				)
			}
		})
	}
}
