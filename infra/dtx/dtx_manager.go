package dtx

import "context"

func NewDistributeTxManager() *DistributeTxManager {
	return &DistributeTxManager{}
}

type DistributeTxManager struct {
	saga     []*Saga
	handlers map[string]Handler
}

func (d *DistributeTxManager) Start(ctx context.Context) error {
	return nil
}

func (d *DistributeTxManager) NewTx(ctx context.Context) *Saga {
	return nil
}

func (d *DistributeTxManager) RegisterHandler(ctx context.Context) error {
	return nil
}

type Handler func(ctx context.Context, payload []byte) error

type HandlerName string

type Action struct {
	Name   HandlerName
	Depend []HandlerName
}
