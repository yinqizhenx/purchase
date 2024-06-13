package app_service

import (
	"purchase/app/assembler"
	"purchase/domain/repo"
	"purchase/domain/service"
)

type SuAppSrv struct {
	suSrv       *service.TicketSupplyDomainSrv
	suRepo      repo.SURepo
	suAssembler *assembler.SuAssembler
}

func NewSuAppSrv(suSrv *service.TicketSupplyDomainSrv, suRepo repo.SURepo) *SuAppSrv {
	return &SuAppSrv{suSrv: suSrv, suRepo: suRepo}
}

// func (s *SuAppSrv) DispatchTicketApply(ctx context.Context, req *pb.AddPAInfo) (*pb.AddPaymentApplyRes, error) {
// 	// 获取参数
// 	// ...
// 	codes := make([]string, 0)
// 	r, err := s.suSrv.NewEmptySuFromPA(ctx, codes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = r.ValidateTasksSevenSame(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return s.suAssembler.SuToDispatchPB(r), nil
// }
