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
	TransID    string
	Name       string
	State      string
	FinishedAt time.Time
	CreatedAt  time.Time
	CreatedBy  string
	UpdatedAt  time.Time
	UpdatedBy  string
}

type Branch struct {
	BranchID         string
	TransID          string
	Type             string
	State            string
	Name             string
	Executor         string
	Payload          string
	ActionDepend     string
	CompensateDepend string
	FinishedAt       time.Time
	IsDead           bool
	CreatedAt        time.Time
	CreatedBy        string
	UpdatedAt        time.Time
	UpdatedBy        string
}

type BranchType string

const (
	SagaAction     BranchType = "saga_action"
	SagaCompensate BranchType = "saga_compensate"
)

type TransStorage interface {
	SaveTrans(ctx context.Context, tx *Trans) error
	SaveBranch(ctx context.Context, branchList []*Branch) error
	UpdateTransState(ctx context.Context, transID, newState string) error
	UpdateBranchState(ctx context.Context, branchID, newState string) error
}
