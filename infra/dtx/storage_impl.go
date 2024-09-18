package dtx

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func (s *StorageImpl) SaveTrans(ctx context.Context, t *Trans) error {
	err := s.getTransClient(ctx).Create().
		SetTransID(t.TransID).
		SetName(t.Name).
		SetState(t.State).
		SetFinishedAt(t.FinishedAt).
		SetCreatedAt(t.CreatedAt).
		SetCreatedBy(t.CreatedBy).
		SetUpdatedBy(t.UpdatedBy).
		SetUpdatedAt(t.UpdatedAt).
		Exec(ctx)
	return err
}

func (s *StorageImpl) SaveBranch(ctx context.Context, branchList []*Branch) error {
	if len(branchList) == 0 {
		return nil
	}
	buildCreate := make([]*ent.BranchCreate, 0)
	for _, b := range branchList {
		create := s.getBranchClient(ctx).Create().
			SetBranchID(b.BranchID).
			SetTransID(b.TransID).
			SetType(b.Type).
			SetState(b.State.String()).
			SetName(b.Name).
			SetAction(b.Action).
			SetCompensate(b.Compensate).
			SetPayload(b.Payload).
			SetActionDepend(strings.Join(b.ActionDepend, ",")).
			SetCompensateDepend(strings.Join(b.CompensateDepend, ",")).
			SetFinishedAt(b.FinishedAt).
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

func (s *StorageImpl) UpdateTransState(ctx context.Context, transID string, newState string) error {
	err := s.getTransClient(ctx).Update().
		SetState(newState).
		Where(trans.TransID(transID)).
		Exec(ctx)
	return err
}

func (s *StorageImpl) UpdateBranchState(ctx context.Context, branchID string, newState string) error {
	err := s.getBranchClient(ctx).Update().
		SetState(newState).
		Where(branch.BranchID(branchID)).
		Exec(ctx)
	return err
}

func (s *StorageImpl) GetPendingTrans(ctx context.Context) (map[string]*Trans, error) {
	transList, err := s.getTransClient(ctx).Query().Where(trans.State("pending")).All(ctx)
	if err != nil {
		return nil, err
	}
	transMap := make(map[string]*Trans)
	for _, t := range transList {
		transMap[t.TransID] = ConvertTrans(t)
	}
	return transMap, nil
}

func (s *StorageImpl) MustGetBranchesByTransIDList(ctx context.Context, transIDList []string) (map[string][]*Branch, error) {
	branchMap := make(map[string][]*Branch)
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
		TransID:    t.TransID,
		Name:       t.Name,
		State:      t.State,
		FinishedAt: t.FinishedAt,
		CreatedAt:  t.CreatedAt,
		CreatedBy:  t.CreatedBy,
	}
}

func ConvertBranch(b *ent.Branch) *Branch {
	return &Branch{
		BranchID:         b.BranchID,
		TransID:          b.TransID,
		Type:             b.Type,
		State:            StepStatus(b.State),
		Name:             b.Name,
		Action:           b.Action,
		Compensate:       b.Compensate,
		Payload:          b.Payload,
		ActionDepend:     strings.Split(b.ActionDepend, ","),
		CompensateDepend: strings.Split(b.CompensateDepend, ","),
		FinishedAt:       b.FinishedAt,
		CreatedAt:        b.CreatedAt,
		CreatedBy:        b.CreatedBy,
		UpdatedAt:        b.UpdatedAt,
		IsDead:           b.IsDead,
	}
}
