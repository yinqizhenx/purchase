package event

import (
	"context"

	"purchase/domain/repo"
)

type EventService struct {
	eventRepo repo.EventRepo
}

func NewEventService(repo repo.EventRepo) *EventService {
	return &EventService{
		eventRepo: repo,
	}
}

func (s *EventService) PubEventAsync(ctx context.Context, events ...Event) error {
	return s.eventRepo.Save(ctx, events...)
}
