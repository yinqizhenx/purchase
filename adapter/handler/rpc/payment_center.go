package rpc

import (
	"context"

	"purchase/app"
	pb "purchase/idl/payment_center"
)

type PaymentCenterHandler struct {
	srv *app.PaymentCenterAppService
}

func (h *PaymentCenterHandler) AddPaymentApply(ctx context.Context, req *pb.AddPAReq) (*pb.AddPARes, error) {
	return h.srv.AddPaymentApply(ctx, req)
}
