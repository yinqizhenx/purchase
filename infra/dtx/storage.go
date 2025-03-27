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
	ID           int
	Name         string
	State        string
	ExecuteState string
	FinishedAt   time.Time
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
}

type Branch struct {
	Code              string
	TransID           int
	Type              string
	State             StepStatus
	Name              string
	Action            string
	Compensate        string
	ActionPayload     []byte
	CompensatePayload []byte
	ActionDepend      []string
	CompensateDepend  []string
	IsDead            bool
	CreatedAt         time.Time
	CreatedBy         string
	UpdatedAt         time.Time
	UpdatedBy         string
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
	UpdateTransExecuteState(ctx context.Context, transID int, newState string) error
	UpdateBranchState(ctx context.Context, code, newState string) error
	GetExecutingTrans(ctx context.Context) (map[int]*Trans, error)
	MustGetBranchesByTransIDList(ctx context.Context, transIDList []int) (map[int][]*Branch, error)
}
