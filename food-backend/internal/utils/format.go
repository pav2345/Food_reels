package utils

import "strconv"

func FormatFixed1(value float64) string {
	return strconv.FormatFloat(value, 'f', 1, 64)
}
