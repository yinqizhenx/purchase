package app

import (
	"context"

	"purchase/app/assembler"
	"purchase/domain/factory"
	"purchase/domain/repo"
	"purchase/domain/service"
	pb "purchase/idl/payment_center"
	"purchase/infra/persistence/tx"
)

type PaymentCenterAppService struct {
	paSrv     *service.PADomainService
	paRepo    repo.PaymentCenterRepo
	assembler *assembler.Assembler
	txm       *tx.TransactionManager
	pcFactory *factory.PCFactory
}

func NewPaymentCenterAppService(paSrv *service.PADomainService, paRepo repo.PaymentCenterRepo, asb *assembler.Assembler, txm *tx.TransactionManager) *PaymentCenterAppService {
	return &PaymentCenterAppService{
		paSrv:     paSrv,
		paRepo:    paRepo,
		assembler: asb,
		txm:       txm,
	}
}

func (s *PaymentCenterAppService) AddPaymentApply(ctx context.Context, req *pb.AddPAReq) (*pb.AddPAResp, error) {
	pa, err := s.pcFactory.BuildPA(ctx, req)
	if err != nil {
		return nil, err
	}
	err = s.txm.Transaction(ctx, func(ctx context.Context) error {
		err = s.paSrv.AddPA(ctx, pa)
		if err != nil {
			return err
		}
		return s.paSrv.PubEvent(ctx, pa.Events()...)
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddPAResp{Code: req.Code}, nil
}

func (s *PaymentCenterAppService) UpdatePaymentApply(ctx context.Context, req *pb.UpdatePAReq) (*pb.UpdatePAResp, error) {
	pa := s.assembler.PAUpdateDtoToDo(req)
	err := s.txm.Transaction(ctx, func(ctx context.Context) error {
		err := s.paSrv.UpdatePA(ctx, pa)
		if err != nil {
			return err
		}
		return s.paSrv.PubEvent(ctx, pa.Events()...)
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdatePAResp{Code: req.Code}, nil
}
