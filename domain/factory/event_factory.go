package factory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"purchase/domain/entity/payment_center"
	"purchase/domain/event"
	"purchase/infra/utils"
)

func NewEventFactory() *EventFactory {
	return &EventFactory{}
}

type EventFactory struct{}

func (f *EventFactory) NewPACreateEvent(ctx context.Context, h *payment_center.PAHead) *event.PACreated {
	return &event.PACreated{
		EventID:   uuid.New().String(),
		PACode:    h.Code,
		CreatedBy: utils.GetCurrentUser(ctx),
		CreatedAt: time.Now(),
	}
}

func (f *EventFactory) NewPAUpdateEvent(ctx context.Context, h *payment_center.PAHead) *event.PACreated {
	return &event.PACreated{
		EventID:   uuid.New().String(),
		PACode:    h.Code,
		CreatedBy: utils.GetCurrentUser(ctx),
		CreatedAt: time.Now(),
	}
}
