package payment_center

import (
	"purchase/domain/entity/company"
	"purchase/domain/entity/department"
	"purchase/domain/entity/supplier"
	"purchase/domain/entity/user"
	"purchase/domain/vo"
)

// PAHead  付款中心-PA单
type PAHead struct {
	ID         int64       `db:"id" json:"id"`
	Code       string      `db:"code" json:"code"`             //  单号
	State      vo.DocState `db:"state" json:"state"`           //  状态
	PayAmount  string      `db:"pay_amount" json:"pay_amount"` //  付款金额
	Applicant  *user.User  `db:"applicant" json:"applicant"`   //  实际需求人
	Department *department.Department
	Currency   string
	IsAdv      bool
	HasInvoice bool
	Company    *company.Company
	Supplier   *supplier.Supplier
	Remark     string
	Rows       []*PARow
	snapshot   *PAHead
}

func (p *PAHead) IsSubmit() bool {
	return p.State == vo.DocStateSubmitted
}

func (p *PAHead) deepCopy() *PAHead {
	return nil
}

func (p *PAHead) Attach() {
	if p.snapshot == nil || p.snapshot.Code == p.Code {
		p.snapshot = p.deepCopy()
	}
}

func (p *PAHead) Detach() {
	if p.snapshot != nil && p.snapshot.Code == p.Code {
		p.snapshot = nil
	}
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
	diff := &PADiff{
		OrderChanged: p.headChanged(),
	}
	return diff
}

func (p *PAHead) headChanged() bool {
	if p.State != p.snapshot.State {
		return true
	}
	if p.PayAmount != p.snapshot.PayAmount {
		return true
	}
	if p.Applicant.Account != p.snapshot.Applicant.Account {
		return true
	}
	if p.Department.GetPathCodeString() != p.snapshot.Department.GetPathCodeString() || p.Department.I18nPathNames != p.snapshot.Department.I18nPathNames {
		return true
	}
	if p.Currency != p.snapshot.Currency {
		return true
	}
	if p.IsAdv != p.snapshot.IsAdv {
		return true
	}
	if p.HasInvoice != p.snapshot.HasInvoice {
		return true
	}
	if p.Company.Code != p.snapshot.Company.Code {
		return true
	}
	if p.Supplier.Code != p.snapshot.Code {
		return true
	}
	if p.Remark != p.snapshot.Remark {
		return true
	}
	return false
}

func (p *PAHead) ChangeProductCnt() error {
	// ...  // 业务逻辑
	// p.raisePACreateEvent() // 生成exam后，调用发布事件的方法
	return nil
}

// func (p *PAHead) AppendEvent(e event.Event) {
// 	p.events = append(p.events, e)
// }

func (p *PAHead) SetSnapshot(h *PAHead) {
	p.snapshot = h
}

// func (p *PAHead) raisePACreateEvent() {
// 	p.RaiseEvent(&event.PACreated{
// 		EventID: p.Code,
// 	})
// }

// func (p *PAHead) Events() []event.Event {
// 	return p.events
// }

func (p *PAHead) Update(h *PAHead) {

}

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
