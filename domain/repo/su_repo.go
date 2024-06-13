package repo

import (
	"context"

	"purchase/domain/entity/su"
)

type SURepo interface {
	NextIdentity() (int64, error)
	Save(context.Context, *su.SU) error                 // 保存一个聚合
	Find(context.Context, string) (*su.SU, error)       // 通过id查找对应的聚合
	FindNonNil(context.Context, string) (*su.SU, error) // 通过id查找对应的聚合，聚合不存在的话返回错误
	Remove(context.Context, *su.SU) error               // 将一个聚合从仓储中删除
}

type SUTaskRepo interface {
	GetSuTasksByPaCodes(context.Context, []string) ([]*su.SUTask, error)
}
