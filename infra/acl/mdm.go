package acl

import (
	"context"

	"purchase/domain/entity/department"
	"purchase/domain/entity/user"
	"purchase/domain/sal"
	"purchase/infra/request"
)

const MDMServiceURL = "/employee/op/searchName"

func NewMDMService(cli *request.HttpClient) sal.MDMService {
	return &MDMServiceImpl{client: cli}
}

type MDMServiceImpl struct {
	client *request.HttpClient
}

func (mdm *MDMServiceImpl) GetUser(ctx context.Context, account string) (*user.User, error) {
	return &user.User{Account: account}, nil
	// req := NewMDMSearchReq(account)
	// resp := NewMDMSearchResp()
	// err := mdm.client.NewRequest(MDMServiceURL, req, resp).Get(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// return resp.toUser(), nil
}

func (mdm *MDMServiceImpl) GetDepartment(ctx context.Context, code string) (*department.Department, error) {
	return &department.Department{Code: code}, nil
}
