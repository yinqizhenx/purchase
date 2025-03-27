package dtx

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/wire"

	"purchase/infra/persistence/dal/db/ent"
	"purchase/infra/persistence/dal/db/ent/branch"
	"purchase/infra/persistence/dal/db/ent/trans"
	"purchase/infra/persistence/tx"
)

var ProviderSet = wire.NewSet(NewTransStorage)

func NewTransStorage(db *ent.Client) TransStorage {
	return &StorageImpl{
		db: db,
	}
}

type StorageImpl struct {
	db *ent.Client
}

// getClient 确保context在事务中时，数据库操作在事务中执行
func (s *StorageImpl) getTransClient(ctx context.Context) *ent.TransClient {
	txCtx, ok := ctx.(*tx.TransactionContext)
	if ok {
		return txCtx.Tx().Trans
	}
	return s.db.Trans
}

func (s *StorageImpl) getBranchClient(ctx context.Context) *ent.BranchClient {
	txCtx, ok := ctx.(*tx.TransactionContext)
	if ok {
		return txCtx.Tx().Branch
	}
	return s.db.Branch
}

func (s *StorageImpl) SaveTrans(ctx context.Context, t *Trans) (int, error) {
	tran, err := s.getTransClient(ctx).Create().
		SetName(t.Name).
		SetState(t.State).
		SetExecuteState(t.ExecuteState).
		SetFinishedAt(t.FinishedAt).
		SetCreatedAt(t.CreatedAt).
		SetCreatedBy(t.CreatedBy).
		SetUpdatedBy(t.UpdatedBy).
		SetUpdatedAt(t.UpdatedAt).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return tran.ID, nil
}

func (s *StorageImpl) SaveBranch(ctx context.Context, branchList []*Branch) error {
	if len(branchList) == 0 {
		return nil
	}
	buildCreate := make([]*ent.BranchCreate, 0)
	for _, b := range branchList {
		create := s.getBranchClient(ctx).Create().
			SetCode(b.Code).
			SetTransID(b.TransID).
			SetType(b.Type).
			SetState(b.State.String()).
			SetName(b.Name).
			SetAction(b.Action).
			SetCompensate(b.Compensate).
			SetActionPayload(string(b.ActionPayload)).
			SetCompensatePayload(string(b.CompensatePayload)).
			SetActionDepend(strings.Join(b.ActionDepend, ",")).
			SetCompensateDepend(strings.Join(b.CompensateDepend, ",")).
			SetCreatedAt(b.CreatedAt).
			SetCreatedBy(b.CreatedBy).
			SetUpdatedBy(b.UpdatedBy).
			SetUpdatedAt(b.UpdatedAt).
			SetIsDead(b.IsDead)
		buildCreate = append(buildCreate, create)

	}
	_, err := s.getBranchClient(ctx).CreateBulk(buildCreate...).Save(ctx)
	return err
}

func (s *StorageImpl) UpdateTransState(ctx context.Context, transID int, newState string) error {
	err := s.getTransClient(ctx).Update().
		SetState(newState).
		Where(trans.ID(transID)).
		Exec(ctx)
	return err
}

func (s *StorageImpl) UpdateTransExecuteStateDone(ctx context.Context, transID int, newState string) error {
	err := s.getTransClient(ctx).Update().
		SetExecuteState(newState).
		SetFinishedAt(time.Now()).
		Where(trans.ID(transID)).
		Exec(ctx)
	return err
}

func (s *StorageImpl) UpdateBranchState(ctx context.Context, code, newState string) error {
	err := s.getBranchClient(ctx).Update().
		SetState(newState).
		Where(branch.Code(code)).
		Exec(ctx)
	return err
}

func (s *StorageImpl) GetExecutingTrans(ctx context.Context) (map[int]*Trans, error) {
	transList, err := s.getTransClient(ctx).Query().Where(trans.ExecuteState("executing")).All(ctx)
	if err != nil {
		return nil, err
	}
	transMap := make(map[int]*Trans)
	for _, t := range transList {
		transMap[t.ID] = ConvertTrans(t)
	}
	return transMap, nil
}

func (s *StorageImpl) MustGetBranchesByTransIDList(ctx context.Context, transIDList []int) (map[int][]*Branch, error) {
	branchMap := make(map[int][]*Branch)
	branchList, err := s.getBranchClient(ctx).Query().Where(branch.TransIDIn(transIDList...)).All(ctx)
	if err != nil {
		return nil, err
	}
	for _, b := range branchList {
		branchMap[b.TransID] = append(branchMap[b.TransID], ConvertBranch(b))
	}
	for _, transID := range transIDList {
		if _, ok := branchMap[transID]; !ok {
			return nil, errors.New(fmt.Sprintf("branch with transID[%s] not found", transID))
		}
	}
	return branchMap, nil
}

func ConvertTrans(t *ent.Trans) *Trans {
	return &Trans{
		ID:           t.ID,
		Name:         t.Name,
		State:        t.State,
		ExecuteState: t.ExecuteState,
		FinishedAt:   t.FinishedAt,
		CreatedAt:    t.CreatedAt,
		CreatedBy:    t.CreatedBy,
	}
}

func ConvertBranch(b *ent.Branch) *Branch {
	actionDepend := make([]string, 0)
	if b.ActionDepend != "" {
		actionDepend = strings.Split(b.ActionDepend, ",")
	}
	return &Branch{
		Code:              b.Code,
		TransID:           b.TransID,
		Type:              b.Type,
		State:             StepStatus(b.State),
		Name:              b.Name,
		Action:            b.Action,
		Compensate:        b.Compensate,
		ActionPayload:     []byte(b.ActionPayload),
		CompensatePayload: []byte(b.CompensatePayload),
		ActionDepend:      actionDepend,
		CompensateDepend:  strings.Split(b.CompensateDepend, ","),
		CreatedAt:         b.CreatedAt,
		CreatedBy:         b.CreatedBy,
		UpdatedAt:         b.UpdatedAt,
		IsDead:            b.IsDead,
	}
}
