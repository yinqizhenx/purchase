package acl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAiSpamChecker)
