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
	eventService *event.EventService
	pcFactory    *factory.PCFactory
	eventFactory *factory.EventFactory
}

func NewPADomainService(repo repo.PaymentCenterRepo, mdm sal.MDMService, eventService *event.EventService, f *factory.PCFactory, ef *factory.EventFactory) *PADomainService {
	return &PADomainService{
		repo:         repo,
		mdm:          mdm,
		eventService: eventService,
		pcFactory:    f,
		eventFactory: ef,
	}
}

func (s *PADomainService) AddPA(ctx context.Context, pa *payment_center.PAHead) error {
	err := s.repo.Save(ctx, pa)
	if err != nil {
		return err
	}
	e := s.eventFactory.NewPACreateEvent(ctx, pa)
	return s.eventService.PubEventAsync(ctx, e)
}

func (s *PADomainService) UpdatePA(ctx context.Context, pa *payment_center.PAHead) error {
	err := s.repo.Save(ctx, pa)
	if err != nil {
		return err
	}
	e := s.eventFactory.NewPAUpdateEvent(ctx, pa)
	return s.eventService.PubEventAsync(ctx, e)
}
