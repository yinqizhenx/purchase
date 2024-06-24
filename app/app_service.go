package app

import (
	"github.com/google/wire"
	// pb "purchase/idl/payment_center"

	pb "purchase/idl/payment_center"
)

type UnImplementedServer struct {
	pb.UnimplementedPaymentCenterServer
}

type Service struct {
	UnImplementedServer
	*PaymentCenterAppService
	*SuAppService
	// ... // 对其他服务的引用等
}

var ProviderSet = wire.NewSet(NewDomainEventAppService, NewPaymentCenterAppService, NewSuAppService)

func NewPurchaseService(pcSrv *PaymentCenterAppService, srv *SuAppService) *Service {
	return &Service{
		PaymentCenterAppService: pcSrv,
		SuAppService:            srv,
	}
}
