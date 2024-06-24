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
	paSrv       *service.PAService
	paRepo      repo.PaymentCenterRepo
	paAssembler *assembler.PaymentAssembler
	txm         *tx.TransactionManager
}

func NewPaymentCenterAppService(paSrv *service.PAService, paRepo repo.PaymentCenterRepo, paAssembler *assembler.PaymentAssembler, txm *tx.TransactionManager) *PaymentCenterAppService {
	return &PaymentCenterAppService{
		paSrv:       paSrv,
		paRepo:      paRepo,
		paAssembler: paAssembler,
		txm:         txm,
	}
}

func (s *PaymentCenterAppService) AddPaymentApply(ctx context.Context, req *pb.AddPAReq) (*pb.AddPARes, error) {
	pa := s.paAssembler.PAHeadDtoToDo(req)
	err := s.txm.Transaction(ctx, func(ctx context.Context) error {
		return s.paSrv.AddPA(ctx, pa)
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddPARes{Code: req.Code}, nil
}
