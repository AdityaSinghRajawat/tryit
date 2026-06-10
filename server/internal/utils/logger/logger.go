// Package logger wraps log/slog with a fixed format and a Redact helper. All
// log output is redacted at the call site — never log a Secret.Reveal() value.
package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger = slog.Logger

func New(level string) *Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl})
	return slog.New(h)
}
