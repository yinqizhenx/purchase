package app_service

import (
	"context"

	"purchase/app/assembler"
	"purchase/domain/repo"
	"purchase/domain/service"
	pb "purchase/idl/payment_center"
	"purchase/infra/persistence/tx"
)

type PaymentCenterAppSrv struct {
	paSrv       *service.PAService
	paRepo      repo.PaymentCenterRepo
	paAssembler *assembler.PaymentAssembler
	txm         *tx.TransactionManager
}

func NewPaymentCenterAppSrv(paSrv *service.PAService, paRepo repo.PaymentCenterRepo, paAssembler *assembler.PaymentAssembler, txm *tx.TransactionManager) *PaymentCenterAppSrv {
	return &PaymentCenterAppSrv{
		paSrv:       paSrv,
		paRepo:      paRepo,
		paAssembler: paAssembler,
		txm:         txm,
	}
}

func (s *PaymentCenterAppSrv) AddPaymentApply(ctx context.Context, req *pb.AddPAReq) (*pb.AddPARes, error) {
	pa := s.paAssembler.PAHeadDtoToDo(req)
	err := s.txm.Transaction(ctx, func(ctx context.Context) error {
		return s.paSrv.AddPA(ctx, pa)
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddPARes{Code: req.Code}, nil
}
