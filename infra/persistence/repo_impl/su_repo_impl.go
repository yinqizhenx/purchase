package repo_impl

import (
	"context"

	"purchase/domain/entity/su"
	"purchase/domain/repo"
	"purchase/infra/mq"
	"purchase/infra/persistence/dal"
)

type SURepository struct {
	data      *dal.SUDal
	publisher mq.Publisher
}

func NewSURepository(data *dal.SUDal, pub mq.Publisher) repo.SURepo {
	return &SURepository{
		// db: cli,
		data:      data,
		publisher: pub,
	}
}

func (r *SURepository) NextIdentity() (int64, error) {
	return 0, nil
}

func (r *SURepository) Save(ctx context.Context, order *su.SU) error {
	// // 发布事件
	// for _, event := range order.Events() {
	// 	err := r.publisher.Publish(ctx, event)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// order.ClearEvents()
	// // 持久化
	// diff := order.DetectChanges()
	// if diff == nil {
	// 	// diff 为空，说明当前不需要追踪变更，采用全量更新的方式
	// 	// orderPO := converter.OrderToPO(order)
	// 	if err := r.data.AddPA(ctx, order); err != nil {
	// 		return err
	// 	}
	// 	// for _, item := range order.Items {
	// 	// 	itemPO := converter.OrderItemToPO(item)
	// 	if err := r.data.UpsertPARows(ctx, order.Rows); err != nil {
	// 		return err
	// 	}
	// 	// }
	// } else {
	// 	// 根据diff，只更新发生了变更的表
	// 	if diff.OrderChanged {
	// 		// orderPO := converter.OrderToPO(order)
	// 		if err := r.data.UpdatePA(ctx, order); err != nil {
	// 			return err
	// 		}
	// 	}
	//
	// 	if err := r.data.SoftDeleteRows(ctx, diff.RemovedRows); err != nil {
	// 		return err
	// 	}
	//
	// 	if err := r.data.AddPARows(ctx, diff.AddedItems); err != nil {
	// 		return err
	// 	}
	//
	// 	if err := r.data.UpdatePARows(ctx, diff.ModifiedItems); err != nil {
	// 		return err
	// 	}
	// 	// for _, item := range diff.RemovedRows {
	// 	// 	itemPO := converter.OrderItemToPO(item)
	// 	// 	if err != r.orderItemDal.SoftDelete(ctx, itemPO); err != nil {
	// 	// 		return err
	// 	// 	}
	// 	// }
	// 	// for _, item := range diff.AddedItems {
	// 	// 	itemPO := converter.OrderItemToPO(item)
	// 	// 	if err != r.orderItemDal.Create(ctx, itemPO); err != nil {
	// 	// 		return err
	// 	// 	}
	// 	// }
	// 	// for _, item := range diff.ModifiedItems {
	// 	// 	itemPO := converter.OrderItemToPO(item)
	// 	// 	if err != r.orderItemDal.Update(ctx, itemPO); err != nil {
	// 	// 		return err
	// 	// 	}
	// 	// }
	// }
	// order.Attach() // 再次调用Attach开始新一轮的追踪
	return nil
}

func (r *SURepository) Find(ctx context.Context, code string) (*su.SU, error) {
	// order := ... // 获取Exam聚合根的流程
	// if order != nil {
	// 	order.Attach()  // 之后调用Attach方法生成Snapshot，开始追踪
	// }
	// return order, nil
	return nil, nil
}

func (r *SURepository) FindNonNil(ctx context.Context, code string) (*su.SU, error) {
	// ...  // 同之前逻辑
	return nil, nil
}

func (r *SURepository) Remove(ctx context.Context, order *su.SU) error {
	// ... // 调用dal执行删除逻辑
	// order.Detach()  // 删除掉 Snapshot 不再追踪
	return nil
}
