package service

import (
	"context"
	"time"

	"purchase/domain/entity/payment_center"
	"purchase/domain/event"
	"purchase/domain/repo"
	"purchase/domain/sal"
)

// // 在 domain.service 中定义领域服务的接口
// type SpamChecker interface {
// 	Check(ctx context.Context, content string) error
// }

type PADomainService struct {
	// repository.CommentRepository 和 ContentSal 在领域层都是以接口的形式存在
	// 因为在领域层不关心具体的实现
	repo      repo.PaymentCenterRepo
	mdm       sal.MDMService
	eventRepo repo.EventRepo
}

func NewPADomainService(repo repo.PaymentCenterRepo, mdm sal.MDMService, eventRepo repo.EventRepo) *PADomainService {
	return &PADomainService{
		repo:      repo,
		mdm:       mdm,
		eventRepo: eventRepo,
	}
}

// PubEvent 插入时间表，后台任务异步发送
func (s *PADomainService) PubEvent(ctx context.Context, events ...event.Event) error {
	return s.eventRepo.Save(ctx, events...)
}

func (s *PADomainService) AddPA(ctx context.Context, pa *payment_center.PAHead) error {
	err := s.repo.Save(ctx, pa)
	if err != nil {
		return err
	}
	paCreated := &event.PACreated{
		EventID:     "11111",
		PACode:      pa.Code,
		AggregateID: 123,
		CreatedBy:   "q",
		CreatedAt:   time.Now(),
	}
	return s.PubEvent(ctx, paCreated)
}
