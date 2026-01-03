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

func MustParseStrToFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}

	return f
}
