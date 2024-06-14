package rpc

import (
	"github.com/google/wire"

	"purchase/adapter/scheduler"
)

var ProviderSet = wire.NewSet(NewGRPCServer, NewHttpServer, scheduler.NewAsyncTaskServer, NewDomainEventServer)
