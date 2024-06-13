package su

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type SU struct {
	code string // 单据编号
	// paCodes      []string // pa 编号
	// companyCode  string   // 支付主体编号
	// companyName  string   // 支付主体名称
	// supplierCode string   // '供应商编号',
	// supplierName string   // '供应商名称',
	// contract     string   // 合同;// '合同编号',
	// currency     string   // 币种
	docBrief     string   // 备注
	invList      []string // 发票附件
	otherAttList []string // 其他附件
	createdBy    string
	department   string
	createdAt    time.Time
	submittedAt  time.Time
	tasks        []*SUTask
	invAmt       decimal.Decimal // 补发票金额
	grnAmt       decimal.Decimal // 补验收金额
}

func NewEmptySuFromTasks(ctx context.Context, tasks []*SUTask) *SU {
	su := &SU{tasks: tasks}
	return su
}

func (s *SU) ValidateTasksSevenSame(ctx context.Context) error {
	m := make(map[[6]string]struct{}, 0)
	for _, v := range s.tasks {
		uk := [6]string{
			v.GetCompanyCode(),
			// v.SupplierCode,
			// v.ContractCode,
			// v.Currency,
			// v.Owner,
			// string(v.Head.OrderType),
			// todo 加上template
		}
		m[uk] = struct{}{}
	}
	if len(m) != 1 {
		return errors.New("ctx, ecode.ErrWrongMutiPayments, s.paCodes")
	}
	return nil
}
