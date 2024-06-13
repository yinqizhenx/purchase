package acl

import (
	"context"

	"purchase/domain/entity/department"
	"purchase/domain/entity/user"
)

// MDMService 请求外部服务的接口，防腐层，具体实现在infra
type MDMService interface {
	GetUsers(ctx context.Context, account string) (*user.User, error)
	GetDepartment(ctx context.Context, code string) (*department.Department, error)
}
