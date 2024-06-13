package middleware

import (
	"context"
	"log/slog"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/google/uuid"

	"purchase/infra/logx"
)

type RequestIDKey struct{}

const RequestIDHeader = "request_id"

func InjectContext() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			newCtx := WithLogRequestID(ctx)
			return handler(newCtx, req)
		}
	}
}

func WithLogRequestID(ctx context.Context) context.Context {
	if reqID, ok := ctx.Value(RequestIDHeader).(string); ok && reqID != "" {
		return logx.AppendCtx(ctx, slog.String(RequestIDHeader, reqID))
	}
	reqID, _ := uuid.NewUUID()
	return logx.AppendCtx(ctx, slog.String(RequestIDHeader, reqID.String()))
}
