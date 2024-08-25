package dtx

import (
	"context"
	"time"
)

func NewTrans() *Trans {
	return &Trans{}
}

func NewBranch(name string) *Branch {
	return &Branch{Name: name}
}

type Trans struct {
	TxID       string
	Name       string
	State      string
	FinishedAt time.Time
	CreatedAt  time.Time
	CreatedBy  string
}

type Branch struct {
	ID               string
	TxID             string
	Type             string
	State            string
	Name             string
	Executor         string
	Payload          []byte
	ActionDepend     string
	CompensateDepend string
	FinishedAt       time.Time
	IsDead           bool
	CreatedAt        time.Time
	CreatedBy        string
}

type TransStorage interface {
	SaveTrans(ctx context.Context, tx *Trans) error
	SaveBranch(ctx context.Context, branchList []*Branch) error
	UpdateTransState(ctx context.Context, tx *Trans, newState string) error
	UpdateBranchState(ctx context.Context, branch *Branch, newState string) error
}
