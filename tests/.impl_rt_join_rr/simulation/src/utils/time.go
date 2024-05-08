package utils

import "time"

func MeasureDuration(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}
