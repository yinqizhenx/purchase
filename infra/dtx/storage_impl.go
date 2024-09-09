package dtx

import (
	"context"

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
			SetState(b.State).
			SetExecutor(b.Executor).
			SetPayload(b.Payload).
			SetActionDepend(b.ActionDepend).
			SetCompensateDepend(b.CompensateDepend).
			SetFinishedAt(b.FinishedAt).
			SetCreatedAt(b.CreatedAt).
			SetCreatedBy(b.CreatedBy).
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
