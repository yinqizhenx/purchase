package app_service

import (
	"github.com/google/wire"

	pb "purchase/idl/payment_center"
)

type UnImplementedServer struct {
	pb.UnimplementedPaymentCenterServer
}

type Service struct {
	UnImplementedServer
	*PaymentCenterAppSrv
	*SuAppSrv
	// ... // 对其他服务的引用等
}

var ProviderSet = wire.NewSet(NewPaymentCenterAppSrv, NewSuAppSrv)

func NewPurchaseService(pcSrv *PaymentCenterAppSrv, srv *SuAppSrv) *Service {
	return &Service{
		PaymentCenterAppSrv: pcSrv,
		SuAppSrv:            srv,
	}
}
