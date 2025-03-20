package logx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// Debug logs a message at debug level.
func Debug(ctx context.Context, msg string, a ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelDebug, msg, a...)
}

// Debugf logs a message at debug level.
func Debugf(ctx context.Context, format string, a ...interface{}) {
	slog.DebugContext(ctx, fmt.Sprintf(format, a...))
}

// Info logs a message at info level.
func Info(ctx context.Context, msg string, a ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelInfo, msg, a...)
}

// Infof logs a message at info level.
func Infof(ctx context.Context, format string, a ...interface{}) {
	slog.InfoContext(ctx, fmt.Sprintf(format, a...))
}

// Warn logs a message at warn level.
func Warn(ctx context.Context, msg string, a ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelWarn, msg, a...)
}

// // Warnf logs a message at warnf level.
func Warnf(ctx context.Context, format string, a ...interface{}) {
	slog.WarnContext(ctx, fmt.Sprintf(format, a...))
}

// Error logs a message at error level.
func Error(ctx context.Context, msg string, a ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelError, msg, a...)
}

// Errorf logs a message at error level.
func Errorf(ctx context.Context, format string, a ...interface{}) {
	slog.ErrorContext(ctx, fmt.Sprintf(format, a...))
}

var ProviderSet = wire.NewSet(NewKratosLogger)

func NewKratosLogger() log.Logger {
	return NewSlogImpl()
}
