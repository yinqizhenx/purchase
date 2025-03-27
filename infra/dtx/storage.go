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
	ID         int
	Name       string
	State      string
	FinishedAt time.Time
	CreatedAt  time.Time
	CreatedBy  string
	UpdatedAt  time.Time
	UpdatedBy  string
}

type Branch struct {
	ID               int
	TransID          int
	Type             string
	State            StepStatus
	Name             string
	Action           string
	Compensate       string
	Payload          string
	ActionDepend     []string
	CompensateDepend []string
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
	SaveTrans(ctx context.Context, tx *Trans) (int, error)
	SaveBranch(ctx context.Context, branchList []*Branch) error
	UpdateTransState(ctx context.Context, transID int, newState string) error
	UpdateBranchState(ctx context.Context, branchID int, newState string) error
	GetPendingTrans(ctx context.Context) (map[int]*Trans, error)
	MustGetBranchesByTransIDList(ctx context.Context, transIDList []int) (map[int][]*Branch, error)
}
