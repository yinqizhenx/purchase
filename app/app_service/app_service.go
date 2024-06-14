package app_service

import (
	"github.com/google/wire"

	pb "purchase/idl/payment_center"
)

type Service struct {
	pb.UnimplementedPaymentCenterServer
	*PaymentCenterAppService
	*SuAppService
	// ... // 对其他服务的引用等
}

var ProviderSet = wire.NewSet(NewPaymentCenterAppService, NewSuAppService)

func NewPurchaseService(pcSrv *PaymentCenterAppService, srv *SuAppService) *Service {
	return &Service{
		PaymentCenterAppService: pcSrv,
		SuAppService:            srv,
	}
}
