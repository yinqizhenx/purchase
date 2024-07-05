package service

import (
	"context"

	"purchase/domain/entity/payment_center"
	"purchase/domain/event"
	"purchase/domain/factory"
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
	repo         repo.PaymentCenterRepo
	mdm          sal.MDMService
	eventRepo    repo.EventRepo
	pcFactory    *factory.PCFactory
	eventFactory *factory.EventFactory
}

func NewPADomainService(repo repo.PaymentCenterRepo, mdm sal.MDMService, eventRepo repo.EventRepo, f *factory.PCFactory, ef *factory.EventFactory) *PADomainService {
	return &PADomainService{
		repo:         repo,
		mdm:          mdm,
		eventRepo:    eventRepo,
		pcFactory:    f,
		eventFactory: ef,
	}
}

// PubEvent 插入时间表，后台任务异步发送
func (s *PADomainService) PubEvent(ctx context.Context, events ...event.Event) error {
	return s.eventRepo.Save(ctx, events...)
}

func (s *PADomainService) AddPA(ctx context.Context, pa *payment_center.PAHead) error {
	// head, err := s.pcFactory.BuildPA(ctx, pa)
	// if err != nil {
	// 	return err
	// }
	e := s.eventFactory.NewPACreateEvent(ctx, pa)
	pa.RaiseEvent(e)
	return s.repo.Save(ctx, pa)
}

func (s *PADomainService) UpdatePA(ctx context.Context, update *payment_center.PAHead) error {
	pa, err := s.repo.Find(ctx, update.Code)
	if err != nil {
		return err
	}
	err = s.pcFactory.UpdatePA(ctx, pa, update)
	if err != nil {
		return err
	}
	e := s.eventFactory.NewPAUpdateEvent(ctx, pa)
	pa.RaiseEvent(e)
	return s.repo.Save(ctx, pa)
}
