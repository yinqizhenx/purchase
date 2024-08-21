package dtx

import "context"

func NewDistributeTxManager() *DistributeTxManager {
	return &DistributeTxManager{}
}

type DistributeTxManager struct {
	steps []Step
	// handler
}

func (d *DistributeTxManager) Exec(ctx context.Context) error {
	return nil
}

func (d *DistributeTxManager) RegisterHandler(ctx context.Context) error {
	return nil
}
