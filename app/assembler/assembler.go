package assembler

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewSuAssembler, NewPaymentAssembler)
