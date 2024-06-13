package dal

import (
	"context"

	"purchase/domain/entity/payment_center"
	"purchase/domain/entity/su"
	"purchase/infra/persistence/dal/db/ent"
)

type SUDal struct {
	db *ent.Client
}

func NewSUDal(cli *ent.Client) *SUDal {
	return &SUDal{
		db: cli,
	}
}

func (dal *SUDal) AddSU(ctx context.Context, pa *su.SU) error {
	// err := dal.db.PAHead.Create().
	// 	SetCode(pa.Code).
	// 	SetState(pa.State).
	// 	SetApplicant(pa.Applicant.Account).
	// 	Exec(ctx)
	return nil
}

func (dal *SUDal) UpdateSU(ctx context.Context, pa *payment_center.PAHead) error {
	return nil
}

func (dal *SUDal) AddSURows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *SUDal) UpdateSURows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *SUDal) UpsertSURows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *SUDal) SoftDeleteRows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *SUDal) GetById(ctx context.Context, id int) (*payment_center.PARow, error) {
	return nil, nil
}

func (dal *SUDal) ListById(ctx context.Context, ids []int) ([]*su.SU, error) {
	return nil, nil
}
