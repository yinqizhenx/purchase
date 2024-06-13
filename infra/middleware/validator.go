package middleware

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"google.golang.org/protobuf/proto"
)

var v *protovalidate.Validator

func init() {
	initValidator()
}

func initValidator() {
	var err error
	v, err = protovalidate.New()
	if err != nil {
		panic(err)
	}
}

// WithValidator 使用protovalidate工具
func WithValidator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			err = v.Validate(req.(proto.Message))
			if err != nil {
				return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
			}
			return handler(ctx, req)
		}
	}
}
