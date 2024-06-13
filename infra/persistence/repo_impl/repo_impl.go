package repo_impl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewPARepository, NewSURepository, NewEventRepository)
