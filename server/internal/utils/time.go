package utils

import "time"

// GetCurrTimeStamp returns the current UTC timestamp. The monotonic clock
// reading is preserved on the result so time.Since stays sleep-safe.
func GetCurrTimeStamp() time.Time {
	return time.Now().UTC()
}

// ComputeTimeTaken returns the milliseconds elapsed since start using
// monotonic time — laptop-sleep gaps are not counted.
func ComputeTimeTaken(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}
