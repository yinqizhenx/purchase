package app

import (
	"context"

	"purchase/app/assembler"
	"purchase/domain/factory"
	"purchase/domain/repo"
	"purchase/domain/service"
	pb "purchase/idl/payment_center"
	"purchase/infra/dtx"
	"purchase/infra/persistence/tx"
)

type PaymentCenterAppService struct {
	paSrv     *service.PADomainService
	paRepo    repo.PaymentCenterRepo
	assembler *assembler.Assembler
	txm       *tx.TransactionManager
	pcFactory *factory.PCFactory
	dtm       *dtx.DistributeTxManager
}

func NewPaymentCenterAppService(paSrv *service.PADomainService, paRepo repo.PaymentCenterRepo, asb *assembler.Assembler, txm *tx.TransactionManager, pcFactory *factory.PCFactory, dtm *dtx.DistributeTxManager) *PaymentCenterAppService {
	return &PaymentCenterAppService{
		paSrv:     paSrv,
		paRepo:    paRepo,
		assembler: asb,
		txm:       txm,
		pcFactory: pcFactory,
		dtm:       dtm,
	}
}

func (s *PaymentCenterAppService) AddPaymentApply(ctx context.Context, req *pb.AddPAReq) (*pb.AddPAResp, error) {
	pa, err := s.pcFactory.BuildPA(ctx, req)
	if err != nil {
		return nil, err
	}

	sg, err := s.dtm.NewTransSaga(ctx, dtx.StepTest)
	if err != nil {
		return nil, err
	}
	err = sg.Exec(ctx)
	if err != nil {
		return nil, err
	}

	err = s.txm.Transaction(ctx, func(ctx context.Context) error {
		return s.paSrv.SavePA(ctx, pa)
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddPAResp{Code: req.Code}, nil
}

func (s *PaymentCenterAppService) UpdatePaymentApply(ctx context.Context, req *pb.UpdatePAReq) (*pb.UpdatePAResp, error) {
	pa, err := s.paRepo.MustFind(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	err = s.pcFactory.UpdateBuildPA(ctx, pa, req)
	if err != nil {
		return nil, err
	}
	err = s.txm.Transaction(ctx, func(ctx context.Context) error {
		return s.paSrv.SavePA(ctx, pa)
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdatePAResp{Code: req.Code}, nil
}
