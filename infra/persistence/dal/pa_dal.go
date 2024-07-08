package dal

import (
	"context"

	"purchase/infra/persistence/convertor"
	"purchase/infra/persistence/dal/db/ent/pahead"
	"purchase/infra/persistence/dal/db/ent/parow"

	"purchase/domain/entity/payment_center"
	"purchase/infra/persistence/dal/db/ent"
	"purchase/infra/persistence/tx"
)

type PADal struct {
	db        *ent.Client
	convertor *convertor.Convertor
}

func NewPADal(cli *ent.Client, c *convertor.Convertor) *PADal {
	return &PADal{
		db:        cli,
		convertor: c,
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

// getClient 确保context在事务中时，数据库操作在事务中执行
func (dal *PADal) getRowClient(ctx context.Context) *ent.PARowClient {
	txCtx, ok := ctx.(*tx.TransactionContext)
	if ok {
		return txCtx.Tx().PARow
	}
	return dal.db.PARow
}

func (dal *PADal) InsertPAHead(ctx context.Context, pa *payment_center.PAHead) error {
	err := dal.getClient(ctx).Create().
		SetCode(pa.Code).
		SetState(pa.State.String()).
		SetApplicant(pa.Applicant.Account).
		SetPayAmount(pa.PayAmount).
		SetDepartmentCode(pa.Department.Code).
		SetHasInvoice(pa.HasInvoice).
		SetIsAdv(pa.IsAdv).
		SetSupplierCode(pa.Supplier.Code).
		SetRemark(pa.Remark).
		Exec(ctx)
	return err
}

func (dal *PADal) UpdatePAHead(ctx context.Context, pa *payment_center.PAHead) error {
	return nil
}

func (dal *PADal) UpdatePARows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *PADal) InsertRows(ctx context.Context, rows []*payment_center.PARow) error {
	if len(rows) == 0 {
		return nil
	}
	buildCreate := make([]*ent.PARowCreate, 0)
	for _, row := range rows {
		create := dal.getRowClient(ctx).Create().
			SetHeadCode(row.HeadCode).
			SetRowCode(row.RowCode).
			SetGrnCount(row.GrnCount).
			SetGrnAmount(row.GrnAmount).
			SetPayAmount(row.PayAmount).
			SetDescription(row.DocDescription)
		buildCreate = append(buildCreate, create)

	}
	_, err := dal.getRowClient(ctx).CreateBulk(buildCreate...).Save(ctx)
	return err
}

func (dal *PADal) UpsertPARows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *PADal) SoftDeleteRows(ctx context.Context, rows []*payment_center.PARow) error {
	return nil
}

func (dal *PADal) GetPaHeadByCode(ctx context.Context, code string) (*payment_center.PAHead, error) {
	pa, err := dal.getClient(ctx).Query().Where(pahead.Code(code)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return dal.convertor.ConvertPAHeadPoToDo(ctx, pa)
}

func (dal *PADal) GetPaRowsByHeadCode(ctx context.Context, headCode string) ([]*payment_center.PARow, error) {
	poRows, err := dal.getRowClient(ctx).Query().Where(parow.HeadCode(headCode)).All(ctx)
	if err != nil {
		return nil, err
	}
	rows := make([]*payment_center.PARow, 0)
	for _, row := range poRows {
		rows = append(rows, dal.convertor.ConvertPARowPoToDo(row))
	}
	return rows, nil
}

func (dal *PADal) ListById(ctx context.Context, ids []int) ([]*payment_center.PARow, error) {
	return nil, nil
}
