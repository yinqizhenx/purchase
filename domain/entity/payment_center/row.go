package payment_center

// PARow  付款中心-PA单-明细
type PARow struct {
	HeadCode       string
	RowCode        string
	GrnCount       int32
	GrnAmount      string
	PayAmount      string
	DocDescription string
}
