package utils

import (
	"context"

	"purchase/infra/logx"
)

func SafeGo(ctx context.Context, fn func()) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				logx.Errorf(ctx, "go rourtine panic captured in SafeGo :%v", p)
			}
		}()
		fn()
	}()
}
