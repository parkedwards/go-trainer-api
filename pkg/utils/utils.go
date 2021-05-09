package utils

import "time"

func IsDuringBusinessHours(timestamp time.Time) bool {
	h := timestamp.Hour()
	w := timestamp.Weekday()
	return h >= 8 && h < 17 && w >= 1 && w <= 5
}
