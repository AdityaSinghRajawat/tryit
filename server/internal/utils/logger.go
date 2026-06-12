package utils

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerInstance *zap.Logger
	once           sync.Once
)

// LoggerKey: context key for a per-request *zap.Logger.
type LoggerKey struct{}

// RequestIDKey: context key for the request ID propagated through handlers.
type RequestIDKey struct{}

func getOrCreateLogger() *zap.Logger {
	once.Do(func() {
		cfg := zap.Config{
			Encoding:         "json",
			Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "time",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "msg",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
		}
		var err error
		loggerInstance, err = cfg.Build()
		if err != nil {
			log.Fatalf("error initializing logger: %v", err)
		}
	})
	return loggerInstance
}

// GetLoggerWithoutCtx returns the global singleton.
func GetLoggerWithoutCtx() *zap.Logger { return getOrCreateLogger() }

// GetLogger pulls a request-scoped logger out of ctx, or returns a no-op when absent.
func GetLogger(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(LoggerKey{}).(*zap.Logger); ok {
		return logger
	}
	return zap.NewNop()
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey{}).(string); ok {
		return id
	}
	return "unknown"
}

func LogInfo(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Info(msg, withRequestID(ctx, fields)...)
}

func LogError(ctx context.Context, err error, fields ...zap.Field) {
	GetLogger(ctx).Error(err.Error(), withRequestID(ctx, fields)...)
}

// LogErrorWithStacktrace: use for panics or unexpected errors where the call site matters.
func LogErrorWithStacktrace(ctx context.Context, err error, fields ...zap.Field) {
	GetLogger(ctx).Error(
		err.Error(),
		append(withRequestID(ctx, fields), zap.Stack("stacktrace"))...,
	)
}

func LogDebug(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Debug(msg, withRequestID(ctx, fields)...)
}

// LogInfoWithoutCtx is the bootstrap fallback — prefer LogInfo when ctx is available.
func LogInfoWithoutCtx(msg string, fields ...zap.Field) {
	GetLoggerWithoutCtx().Info(msg, fields...)
}

// LogErrorWithoutCtx is the bootstrap fallback — prefer LogError when ctx is available.
func LogErrorWithoutCtx(err error, fields ...zap.Field) {
	GetLoggerWithoutCtx().Error(err.Error(), fields...)
}

// LogInfoWTimeTaken records "<msg> - Time taken: <ms> ms".
func LogInfoWTimeTaken(ctx context.Context, msg string, start time.Time) {
	LogInfo(ctx, fmt.Sprintf("%s - Time taken: %d ms", msg, ComputeTimeTaken(start)))
}

func withRequestID(ctx context.Context, fields []zap.Field) []zap.Field {
	if id := GetRequestID(ctx); id != "" {
		return append(fields, zap.String("requestId", id))
	}
	return fields
}
