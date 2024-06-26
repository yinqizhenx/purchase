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
	// 发布事件
	// for _, event := range order.Events() {
	// 	m, err := event.Encode()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = r.publisher.Publish(ctx, event.EventName(), m)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// 发布消息需要与业务操作再一个事务里时，写入消息表， todo 直接发布事件与写消息表2选一
	// taskList, err := order.EventTasks()
	// if err != nil {
	// 	return err
	// }
	// if err := r.eventDal.BatchAddTask(ctx, taskList...); err != nil {
	// 	return err
	// }
	//
	// // 在事务提交后，再通知ch
	// tx.RunAfterTxCommit(ctx, func(ctx context.Context) error {
	// 	for _, task := range taskList {
	// 		r.ch.In <- task.TaskID
	// 	}
	// 	return nil
	// })
	//
	// order.ClearEvents()
	// 持久化
	diff := order.DetectChanges()
	if diff == nil {
		// diff 为空，说明当前不需要追踪变更，采用全量更新的方式
		// orderPO := converter.OrderToPO(order)
		if err := r.dal.InsertPA(ctx, order); err != nil {
			return err
		}
		// for _, item := range order.Items {
		// 	itemPO := converter.OrderItemToPO(item)
		if err := r.dal.InsertRows(ctx, order.Rows); err != nil {
			return err
		}
		// }
	} else {
		// 根据diff，只更新发生了变更的表
		if diff.OrderChanged {
			// orderPO := converter.OrderToPO(order)
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
		// for _, item := range diff.RemovedRows {
		// 	itemPO := converter.OrderItemToPO(item)
		// 	if err != r.orderItemDal.SoftDelete(ctx, itemPO); err != nil {
		// 		return err
		// 	}
		// }
		// for _, item := range diff.AddedItems {
		// 	itemPO := converter.OrderItemToPO(item)
		// 	if err != r.orderItemDal.Create(ctx, itemPO); err != nil {
		// 		return err
		// 	}
		// }
		// for _, item := range diff.ModifiedItems {
		// 	itemPO := converter.OrderItemToPO(item)
		// 	if err != r.orderItemDal.Update(ctx, itemPO); err != nil {
		// 		return err
		// 	}
		// }
	}
	order.Attach() // 再次调用Attach开始新一轮的追踪
	return nil
}

func (r *PARepository) Find(ctx context.Context, code string) (*payment_center.PAHead, error) {
	// order := ... // 获取Exam聚合根的流程
	// if order != nil {
	// 	order.Attach()  // 之后调用Attach方法生成Snapshot，开始追踪
	// }
	// return order, nil
	return nil, nil
}

func (r *PARepository) FindNonNil(ctx context.Context, code string) (*payment_center.PAHead, error) {
	// ...  // 同之前逻辑
	return nil, nil
}

func (r *PARepository) Remove(ctx context.Context, order *payment_center.PAHead) error {
	// ... // 调用dal执行删除逻辑
	// order.Detach()  // 删除掉 Snapshot 不再追踪
	return nil
}
