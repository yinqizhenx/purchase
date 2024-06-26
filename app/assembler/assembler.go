package assembler

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAssembler)

type Assembler struct{}

func NewAssembler() *Assembler {
	return &Assembler{}
}
