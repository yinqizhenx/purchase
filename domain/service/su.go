package service

import (
	"context"

	"purchase/domain/entity/su"
	"purchase/domain/repo"
)

type TicketSupplyDomainSrv struct {
	// repository.CommentRepository 和 ContentSal 在领域层都是以接口的形式存在
	// 因为在领域层不关心具体的实现
	suRepo     repo.SURepo
	tm         repo.TMRepo
	suTaskRepo repo.SUTaskRepo
}

func NewTicketSupplyDomainSrv(suRepo repo.SURepo) *TicketSupplyDomainSrv {
	return &TicketSupplyDomainSrv{suRepo: suRepo}
}

func (ts *TicketSupplyDomainSrv) GetInvalidTasks(ctx context.Context, codes []string) error {
	taskStates, err := ts.tm.GetTasksStats(ctx, "", codes)
	if err != nil {
		return err
	}
	for _, v := range taskStates {
		// 失效原因-已完成
		if !v.OperateAble() {
			continue
		}
		// 失效原因-已转交
		if !v.UserHadPermission() {
			return nil
		}
	}
	return nil
}

func (ts *TicketSupplyDomainSrv) GetSuTasksByPaCodes(ctx context.Context, paCode string) ([]*su.SUTask, error) {
	return nil, nil
}

func (ts *TicketSupplyDomainSrv) NewEmptySuFromPA(ctx context.Context, paCodes []string) (*su.SU, error) {
	tasks, err := ts.suTaskRepo.GetSuTasksByPaCodes(ctx, paCodes)
	if err != nil {
		return nil, err
	}
	s := su.NewEmptySuFromTasks(ctx, tasks)
	return s, nil
}
