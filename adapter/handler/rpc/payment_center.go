package rpc

import (
	"context"

	"purchase/app"
	pb "purchase/idl/payment_center"
)

type PaymentCenterHandler struct {
	srv *app.PaymentCenterAppService
}

func (h *PaymentCenterHandler) AddPaymentApply(ctx context.Context, req *pb.AddPAReq) (*pb.AddPAResp, error) {
	return h.srv.AddPaymentApply(ctx, req)
}

func (h *PaymentCenterHandler) UpdatePaymentApply(ctx context.Context, req *pb.UpdatePAReq) (*pb.UpdatePAResp, error) {
	return h.srv.UpdatePaymentApply(ctx, req)
}
