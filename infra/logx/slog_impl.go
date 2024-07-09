package logx

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-kratos/kratos/v2/log"
)

func NewSlogImpl() *SlogImpl {
	return &SlogImpl{}
}

type SlogImpl struct{}

func (sl *SlogImpl) Log(level log.Level, kvs ...interface{}) error {
	kvLen := len(kvs)
	data := make([]slog.Attr, 0, kvLen/2)
	if kvLen > 1 {
		for i := 0; i < kvLen; i += 2 {
			data = append(data, slog.Any(fmt.Sprint(kvs[i]), kvs[i+1]))
		}
	}

	msg := ""
	switch level {
	case log.LevelDebug:
		Debug(context.Background(), msg, data...)
	case log.LevelInfo:
		Info(context.Background(), msg, data...)
	case log.LevelWarn:
		Warn(context.Background(), msg, data...)
	case log.LevelError, log.LevelFatal:
		Error(context.Background(), msg, data...)
	}
	return nil
}
