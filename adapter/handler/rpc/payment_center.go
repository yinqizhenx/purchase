package rpc

import (
	"context"

	"purchase/app"
	pb "purchase/idl/payment_center"
)

type PaymentCenterHandler struct {
	srv *app.PaymentCenterAppService
}

func (h *PaymentCenterHandler) AddOrUpdatePaymentApply(ctx context.Context, req *pb.AddOrUpdatePAReq) (*pb.AddOrUpdatePAResp, error) {
	return h.srv.AddOrUpdatePaymentApply(ctx, req)
}
