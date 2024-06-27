package convertor

import (
	"github.com/google/wire"
	"purchase/domain/sal"
)

var ProviderSet = wire.NewSet(NewConvertor)

type Convertor struct {
	mdm sal.MDMService
}

func NewConvertor(mdm sal.MDMService) *Convertor {
	return &Convertor{
		mdm: mdm,
	}
}
