package app

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDomainEventAppService, NewPaymentCenterAppService, NewSuAppService)
