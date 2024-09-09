package server

import (
	"github.com/google/wire"

	"purchase/adapter/scheduler"
	"purchase/infra/dtx"
)

var ProviderSet = wire.NewSet(NewGRPCServer, NewHttpServer, scheduler.NewAsyncTaskServer, NewEventConsumerServer, dtx.NewDistributeTxManager)
