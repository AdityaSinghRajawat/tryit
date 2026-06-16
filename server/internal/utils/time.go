package utils

import "time"

func GetCurrTimeStamp() time.Time {
	return time.Now().UTC()
}

func ComputeTimeTaken(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}
