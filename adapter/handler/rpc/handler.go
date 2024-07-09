package rpc

import (
	"github.com/google/wire"
	// pb "purchase/idl/payment_center"

	pb "purchase/idl/payment_center"
)

type UnImplementedServer struct {
	pb.UnimplementedPaymentCenterServer
}

var _ pb.PaymentCenterServer = (*Server)(nil)

type Server struct {
	*PaymentCenterHandler
	// ... // 对其他服务的引用等
	UnImplementedServer
}

var ProviderSet = wire.NewSet(NewPurchaseServer, NewPaymentCenterHandler)

func NewPurchaseServer(pcHandler *PaymentCenterHandler) *Server {
	return &Server{
		PaymentCenterHandler: pcHandler,
	}
}
