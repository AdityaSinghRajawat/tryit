package utils

import "time"

func GetCurrTimeStamp() time.Time {
	return time.Now().UTC()
}

// ComputeTimeTaken: milliseconds elapsed since start using monotonic time
// (sleep-safe).
func ComputeTimeTaken(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}
