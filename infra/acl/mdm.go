package acl

import (
	"context"

	"purchase/domain/entity/department"
	"purchase/domain/entity/user"
	"purchase/infra/request"
)

const MDMServiceURL = "/employee/op/searchName"

type MDMServiceImpl struct {
	client *request.HttpClient
}

func (mdm *MDMServiceImpl) GetUsers(ctx context.Context, account string) (*user.User, error) {
	req := NewMDMSearchReq(account)
	resp := NewMDMSearchResp()
	err := mdm.client.NewRequest(MDMServiceURL, req, resp).Get(ctx)
	if err != nil {
		return nil, err
	}
	return resp.toUser(), nil
}

func (mdm *MDMServiceImpl) GetDepartment(ctx context.Context, code string) (*department.Department, error) {
	return nil, nil
}
