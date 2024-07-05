package service

import (
	"context"

	"purchase/domain/entity/payment_center"
	"purchase/domain/event"
	"purchase/domain/factory"
	"purchase/domain/repo"
	"purchase/domain/sal"
)

type PADomainService struct {
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

func (s *PADomainService) SavePA(ctx context.Context, pa *payment_center.PAHead) error {
	err := s.repo.Save(ctx, pa)
	if err != nil {
		return err
	}
	if pa.IsSubmit() {
		e := s.eventFactory.NewPACreateEvent(ctx, pa)
		return s.eventService.PubEventAsync(ctx, e)
	}
	return nil
}
