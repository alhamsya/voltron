package util

import (
	"strconv"
	"strings"
)

func ParseStrToFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	return strconv.ParseFloat(s, 64)
}
