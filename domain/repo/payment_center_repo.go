package repo

import (
	"context"

	"purchase/domain/entity/payment_center"
)

type PaymentCenterRepo interface {
	NextIdentity() (int64, error)
	Save(context.Context, *payment_center.PAHead) error               // 保存一个聚合
	Find(context.Context, string) (*payment_center.PAHead, error)     // 通过id查找对应的聚合
	MustFind(context.Context, string) (*payment_center.PAHead, error) // 通过id查找对应的聚合，聚合不存在的话返回错误
	Remove(context.Context, *payment_center.PAHead) error             // 将一个聚合从仓储中删除
}
