package logx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

//	var defaultLogger struct {
//		lock sync.Mutex
//		log.Logger
//	}
//
//	func SetLogger(lg log.Logger) {
//		defaultLogger.lock.Lock()
//		defer defaultLogger.lock.Unlock()
//		defaultLogger.Logger = lg
//	}
//
// Debug logs a message at debug level.
func Debug(ctx context.Context, msg string, a ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelDebug, msg, a...)
}

// Debugf logs a message at debug level.
func Debugf(ctx context.Context, format string, a ...interface{}) {
	slog.DebugContext(ctx, format, fmt.Sprintf(format, a...))
}

// Info logs a message at info level.
func Info(ctx context.Context, msg string, a ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelInfo, msg, a...)
}

// Infof logs a message at info level.
func Infof(ctx context.Context, format string, a ...interface{}) {
	slog.InfoContext(ctx, format, fmt.Sprintf(format, a...))
}

// Warn logs a message at warn level.
func Warn(ctx context.Context, msg string, a ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelWarn, msg, a...)
}

// // Warnf logs a message at warnf level.
func Warnf(ctx context.Context, format string, a ...interface{}) {
	slog.WarnContext(ctx, format, fmt.Sprintf(format, a...))
}

// Error logs a message at error level.
func Error(ctx context.Context, msg string, a ...slog.Attr) {
	slog.LogAttrs(ctx, slog.LevelError, msg, a...)
}

// Errorf logs a message at error level.
func Errorf(ctx context.Context, format string, a ...interface{}) {
	slog.ErrorContext(ctx, fmt.Sprintf(format, a...))
}

var ProviderSet = wire.NewSet(NewLogger, NewLogHelper)

func NewLogger(cfg config.Config) log.Logger {
	zapLog := newZapLog(cfg)
	var lg log.Logger
	lg = &ZapLogger{logger: zapLog}
	// lg = log.With(lg, "request_id", requestID())
	return lg
}

func NewLogHelper(lg log.Logger) *log.Helper {
	return log.NewHelper(lg)
}

//
// func requestID() log.Valuer {
// 	return func(ctx context.Context) interface{} {
// 		return ctx.Value(middleware.RequestIDKey{})
// 	}
// }

// func insert(arr []interface{}, i int, val interface{}) []interface{} {
// 	arr = append(arr, val)
// 	copy(arr[i+1:], arr[i:])
// 	arr[i] = val
// 	return arr
// }
