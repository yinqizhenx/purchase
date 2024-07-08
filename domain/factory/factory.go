package factory

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewEventFactory, NewPCFactory)
