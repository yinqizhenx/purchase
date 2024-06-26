package dal

import (
	"context"

	"purchase/domain/entity/payment_center"
	"purchase/infra/persistence/dal/db/ent"
	"purchase/infra/persistence/tx"
)

type PADal struct {
	db *ent.Client
}

func NewPADal(cli *ent.Client) *PADal {
	return &PADal{
		db: cli,
	}
}

// getClient 确保context在事务中时，数据库操作在事务中执行
func (dal *PADal) getClient(ctx context.Context) *ent.PAHeadClient {
	txCtx, ok := ctx.(*tx.TransactionContext)
	if ok {
		return txCtx.Tx().PAHead
	}
	return dal.db.PAHead
}

func (dal *PADal) InsertPA(ctx context.Context, pa *payment_center.PAHead) error {
	err := dal.getClient(ctx).Create().
		SetCode(pa.Code).
		SetState(pa.State).
		SetApplicant(pa.Applicant.Account).
		SetPayAmount(pa.PayAmount).
		Exec(ctx)
	return err
}

func (dal *PADal) UpdatePA(ctx context.Context, pa *payment_center.PAHead) error {
	return nil
}

func (dal *PADal) UpdatePARows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *PADal) InsertRows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *PADal) UpsertPARows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *PADal) SoftDeleteRows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *PADal) GetById(ctx context.Context, id int) (*payment_center.PARow, error) {
	return nil, nil
}

func (dal *PADal) ListById(ctx context.Context, ids []int) ([]*payment_center.PARow, error) {
	return nil, nil
}
