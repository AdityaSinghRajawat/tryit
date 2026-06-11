package config

import "time"

type executeConsts struct {
	maxRequestBodyBytes int64
	maxRedirects        int
	shutdownDrain       time.Duration
}

var executeI = &executeConsts{
	maxRequestBodyBytes: 5 * 1024 * 1024,
	maxRedirects:        5,
	shutdownDrain:       5 * time.Second,
}

func GetMaxRequestBodyBytes() int64   { return executeI.maxRequestBodyBytes }
func GetMaxRedirects() int            { return executeI.maxRedirects }
func GetShutdownDrain() time.Duration { return executeI.shutdownDrain }
