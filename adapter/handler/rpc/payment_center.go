package rpc

import (
	"context"

	pb "purchase/idl/payment_center"
)

type PaymentCenterHandler struct{}

func (h *PaymentCenterHandler) AddPaymentApply(ctx context.Context, req *pb.AddPAReq) (*pb.AddPARes, error) {
	return nil, nil
}
