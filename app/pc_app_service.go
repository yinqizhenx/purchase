package app

import (
	"context"

	"purchase/app/assembler"
	"purchase/domain/repo"
	"purchase/domain/service"
	pb "purchase/idl/payment_center"
	"purchase/infra/persistence/tx"
)

type PaymentCenterAppService struct {
	paSrv  *service.PADomainService
	paRepo repo.PaymentCenterRepo
	asb    *assembler.Assembler
	txm    *tx.TransactionManager
}

func NewPaymentCenterAppService(paSrv *service.PADomainService, paRepo repo.PaymentCenterRepo, asb *assembler.Assembler, txm *tx.TransactionManager) *PaymentCenterAppService {
	return &PaymentCenterAppService{
		paSrv:  paSrv,
		paRepo: paRepo,
		asb:    asb,
		txm:    txm,
	}
}

func (s *PaymentCenterAppService) AddOrUpdatePaymentApply(ctx context.Context, req *pb.AddOrUpdatePAReq) (*pb.AddOrUpdatePAResp, error) {
	pa := s.asb.PAHeadDtoToDo(req)
	err := s.txm.Transaction(ctx, func(ctx context.Context) error {
		return s.paSrv.AddOrUpdatePA(ctx, pa)
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddOrUpdatePAResp{Code: req.Code}, nil
}
