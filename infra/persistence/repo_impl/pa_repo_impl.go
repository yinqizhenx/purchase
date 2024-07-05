package repo_impl

import (
	"context"

	"purchase/domain/entity/payment_center"
	"purchase/domain/repo"
	"purchase/infra/persistence/dal"
	"purchase/pkg/chanx"
)

type PARepository struct {
	dal      *dal.PADal
	eventDal *dal.AsyncTaskDal
	ch       *chanx.UnboundedChan[string]
}

func NewPARepository(dal *dal.PADal, ch *chanx.UnboundedChan[string], eventDal *dal.AsyncTaskDal) repo.PaymentCenterRepo {
	return &PARepository{
		// db: cli,
		dal: dal,
		// publisher: pub,
		ch:       ch,
		eventDal: eventDal,
	}
}

func (r *PARepository) NextIdentity() (int64, error) {
	return 0, nil
}

func (r *PARepository) Save(ctx context.Context, order *payment_center.PAHead) error {
	diff := order.DetectChanges()
	if diff == nil {
		if err := r.dal.InsertPA(ctx, order); err != nil {
			return err
		}
		if err := r.dal.InsertRows(ctx, order.Rows); err != nil {
			return err
		}
		// }
	} else {
		// 根据diff，只更新发生了变更的表
		if diff.OrderChanged {
			if err := r.dal.UpdatePA(ctx, order); err != nil {
				return err
			}
		}

		if err := r.dal.SoftDeleteRows(ctx, diff.RemovedRows); err != nil {
			return err
		}

		if err := r.dal.InsertRows(ctx, diff.AddedItems); err != nil {
			return err
		}

		if err := r.dal.UpdatePARows(ctx, diff.ModifiedItems); err != nil {
			return err
		}
	}
	order.Attach() // 再次调用Attach开始新一轮的追踪
	return nil
}

func (r *PARepository) Find(ctx context.Context, code string) (*payment_center.PAHead, error) {
	pa, err := r.dal.GetPaHeadByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	rows, err := r.dal.GetPaRowsByHeadCode(ctx, code)
	if err != nil {
		return nil, err
	}
	pa.Rows = rows
	pa.Attach() // 之后调用Attach方法生成Snapshot，开始追踪
	return pa, nil
}

func (r *PARepository) MustFind(ctx context.Context, code string) (*payment_center.PAHead, error) {
	pa, err := r.dal.GetPaHeadByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	rows, err := r.dal.GetPaRowsByHeadCode(ctx, code)
	if err != nil {
		return nil, err
	}
	pa.Rows = rows
	pa.Attach() // 之后调用Attach方法生成Snapshot，开始追踪
	return pa, nil
}

func (r *PARepository) Remove(ctx context.Context, order *payment_center.PAHead) error {
	// ... // 调用dal执行删除逻辑
	// order.Detach()  // 删除掉 Snapshot 不再追踪
	return nil
}
