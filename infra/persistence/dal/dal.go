package dal

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewPADal, NewSUDal, NewAsyncTaskDal, NewFailedMessageDal)
