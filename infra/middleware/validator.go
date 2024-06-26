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

// func (v *Validate) validateUnaryServerInterceptor() grpc.UnaryServerInterceptor {
// 	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 		switch req.(type) {
// 		case proto.Message:
// 			if err := v.v.Validate(req.(proto.Message)); err != nil {
// 				var valErr *protovalidate.ValidationError
// 				if ok := errors.As(err, &valErr); ok && len(valErr.ToProto().GetViolations()) > 0 {
// 					return nil, status.Error(codes.InvalidArgument, valErr.ToProto().GetViolations()[0].Message)
// 				}
// 				return nil, status.Error(codes.InvalidArgument, err.Error())
// 			}
// 		}
// 		return handler(ctx, req)
// 	}
// }
