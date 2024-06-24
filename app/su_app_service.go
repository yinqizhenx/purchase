package app

import (
	"purchase/app/assembler"
	"purchase/domain/repo"
	"purchase/domain/service"
)

type SuAppService struct {
	suSrv       *service.TicketSupplyDomainSrv
	suRepo      repo.SURepo
	suAssembler *assembler.SuAssembler
}

func NewSuAppService(suSrv *service.TicketSupplyDomainSrv, suRepo repo.SURepo) *SuAppService {
	return &SuAppService{suSrv: suSrv, suRepo: suRepo}
}

// func (s *SuAppService) DispatchTicketApply(ctx context.Context, req *pb.AddPAInfo) (*pb.AddPaymentApplyRes, error) {
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
