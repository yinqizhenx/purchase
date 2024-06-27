package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"purchase/domain/entity/payment_center"
	"purchase/domain/event"
	"purchase/domain/factory"
	"purchase/domain/repo"
	"purchase/domain/sal"
	"purchase/infra/utils"
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
	factory   *factory.PCFactory
}

func NewPADomainService(repo repo.PaymentCenterRepo, mdm sal.MDMService, eventRepo repo.EventRepo, f *factory.PCFactory) *PADomainService {
	return &PADomainService{
		repo:      repo,
		mdm:       mdm,
		eventRepo: eventRepo,
		factory:   f,
	}
}

// PubEvent 插入时间表，后台任务异步发送
func (s *PADomainService) PubEvent(ctx context.Context, events ...event.Event) error {
	return s.eventRepo.Save(ctx, events...)
}

func (s *PADomainService) AddOrUpdatePA(ctx context.Context, pa *payment_center.PAHead) error {
	if pa.Code == "" {
		return s.AddPA(ctx, pa)
	}
	return s.UpdatePA(ctx, pa)
}

func (s *PADomainService) AddPA(ctx context.Context, pa *payment_center.PAHead) error {
	head, err := s.factory.BuildPA(ctx, pa)
	if err != nil {
		return err
	}
	head.AppendEvent(s.buildPACreateEvent(ctx, head))
	return s.repo.Save(ctx, head)
}

func (s *PADomainService) UpdatePA(ctx context.Context, update *payment_center.PAHead) error {
	pa, err := s.repo.Find(ctx, update.Code)
	if err != nil {
		return err
	}
	err = s.factory.UpdatePA(ctx, pa, update)
	if err != nil {
		return err
	}
	pa.AppendEvent(s.buildPAUpdateEvent(ctx, pa))
	return s.repo.Save(ctx, pa)
}

func (s *PADomainService) buildPACreateEvent(ctx context.Context, h *payment_center.PAHead) event.Event {
	return &event.PACreated{
		EventID:   uuid.New().String(),
		PACode:    h.Code,
		CreatedBy: utils.GetCurrentUser(ctx),
		CreatedAt: time.Now(),
	}
}

func (s *PADomainService) buildPAUpdateEvent(ctx context.Context, h *payment_center.PAHead) event.Event {
	return &event.PACreated{
		EventID:   uuid.New().String(),
		PACode:    h.Code,
		CreatedBy: utils.GetCurrentUser(ctx),
		CreatedAt: time.Now(),
	}
}
