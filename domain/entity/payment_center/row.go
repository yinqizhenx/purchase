package payment_center

// PARow  付款中心-PA单-明细
type PARow struct {
	ID            int64  `db:"id" json:"id"`
	OrderCode     string `db:"order_code" json:"order_code"`
	DocCode       string `db:"doc_code" json:"doc_code"`             //  单据编号
	RowCurrency   string `db:"row_currency" json:"row_currency"`     //  付款申请任务行原始币种，不是预付单/结算单的币种
	TaxRatio      string `db:"tax_ratio" json:"tax_ratio"`           //  税率
	InitialAmount string `db:"initial_amount" json:"initial_amount"` // 单据金额(不含超额)
}
