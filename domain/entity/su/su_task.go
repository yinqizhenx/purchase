package su

import (
	"time"
)

type SUTask struct {
	code          string   // 单据编号
	paCodes       []string // pa 编号
	companyCode   string   // 支付主体编号
	companyName   string   // 支付主体名称
	supplierCode  string   // '供应商编号',
	supplierName  string   // '供应商名称',
	contract      string   // 合同;// '合同编号',
	currency      string   // 币种
	orderType     string
	paymentAmount string
	applicant     string
	createdBy     string
	department    string
	createdAt     time.Time
	// needInvAmt   decimal.Decimal // 补发票金额
	// needGrnAmt   decimal.Decimal // 补验收金额
	tmCode string // 对tm task引用
}

func (s SUTask) GetCompanyCode() string {
	return s.companyCode
}
