package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
)

const CallbackBeforeResponseKey = "callbackBeforeResponseKey"

func WithCallBackBeforeResponse() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			resp, err := handler(ctx, req)
			if err != nil {
				return nil, err
			}
			v := ctx.Value(CallbackBeforeResponseKey)
			if v != nil {
				if fn, ok := v.(func(context.Context) error); ok {
					err = fn(ctx)
					if err != nil {
						return nil, err
					}
				}
			}
			return resp, nil
		}
	}
}

func InjectCallBackBeforeResponse(ctx context.Context, fn func(context.Context) error) context.Context {
	return context.WithValue(ctx, CallbackBeforeResponseKey, fn)
}
