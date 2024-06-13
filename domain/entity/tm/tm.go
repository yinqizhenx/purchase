package tm

type TMState struct{}

func (s TMState) OperateAble() bool {
	return true
}

func (s TMState) UserHadPermission() bool {
	return true
}
