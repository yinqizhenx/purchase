package app

import (
	"github.com/google/wire"

	"purchase/app/event_handler"
)

var ProviderSet = wire.NewSet(event_handler.NewDomainEventAppService, NewPaymentCenterAppService, NewSuAppService)
