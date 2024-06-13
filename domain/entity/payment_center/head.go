package payment_center

import (
	"purchase/domain/entity/user"
)

// PAHead  付款中心-PA单
type PAHead struct {
	ID        int64     `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`             //  单号
	State     string    `db:"state" json:"state"`           //  状态
	PayAmount string    `db:"pay_amount" json:"pay_amount"` //  付款金额
	Applicant user.User `db:"applicant" json:"applicant"`   //  实际需求人
	Rows      []*PARow
	// events    []event.PAEvent
	snapshot *PAHead
}

func deepCopy(p PAHead) *PAHead {
	return nil
}

func (p *PAHead) Attach() {
	// if e.snapshot == nil || e.snapshot.Id == e.Id {
	// 	e.snapshot = deepCopy(e)
	// }
}

func (p *PAHead) Detach() {
	// if e.snapshot != nil && e.snapshot.Id == e.Id {
	// 	e.snapshot = nil
	// }
}

type PADiff struct {
	OrderChanged  bool
	RemovedRows   []*PARow
	AddedItems    []*PARow
	ModifiedItems []*PARow
}

func (p *PAHead) DetectChanges() *PADiff {
	if p.snapshot == nil {
		return nil
	}
	// 其他diff逻辑
	return nil
}

func (p *PAHead) ChangeProductCnt() error {
	// ...  // 业务逻辑
	// p.raisePACreateEvent() // 生成exam后，调用发布事件的方法
	return nil
}

// func (p *PAHead) raisePACreateEvent() {
// 	p.RaiseEvent(&event.PACreated{
// 		EventID: p.Code,
// 	})
// }

// func (p *PAHead) Events() []event.PAEvent {
// 	return p.events
// }
//
// func (p *PAHead) EventTasks() ([]*async_task.AsyncTask, error) {
// 	msgList := make([]*async_task.AsyncTask, 0, len(p.events))
// 	for _, e := range p.events {
// 		m, err := e.ToAsyncTask()
// 		if err != nil {
// 			return nil, err
// 		}
// 		msgList = append(msgList, m)
// 	}
// 	return msgList, nil
// }
//
// func (p *PAHead) RaiseEvent(event event.PAEvent) {
// 	p.events = append(p.events, event)
// }
//
// func (p *PAHead) ClearEvents() {
// 	for idx := range p.events {
// 		p.events[idx] = nil
// 	}
// 	p.events = nil
// }

// Validate 校验有效性
func (p *PAHead) Validate() error {
	return nil
}

// // 同包下定义校验器
// type PAValidator struct {
// 	order *Order
// 	...
// }
//
// func (v *OrderValidator) Validate() error {
// }
