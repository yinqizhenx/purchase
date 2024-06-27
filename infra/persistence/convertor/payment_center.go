package convertor

import (
	"context"
	"purchase/domain/entity/payment_center"
	"purchase/domain/vo"
	"purchase/infra/persistence/dal/db/ent"
)

func (c *Convertor) ConvertPAHeadPoToDo(ctx context.Context, h *ent.PAHead) (*payment_center.PAHead, error) {
	applicant, err := c.mdm.GetUser(ctx, h.Applicant)
	if err != nil {
		return nil, err
	}

	dept, err := c.mdm.GetDepartment(ctx, h.DepartmentCode)
	if err != nil {
		return nil, err
	}
	return &payment_center.PAHead{
		ID:         h.ID,
		Code:       h.Code,
		State:      vo.DocState(h.State),
		PayAmount:  h.PayAmount,
		Applicant:  applicant,
		Department: dept,
		//Currency: h.
		IsAdv:      h.IsAdv,
		HasInvoice: h.HasInvoice,
		//Company    *company.Company
		//Supplier   *supplier.Supplier
		Remark: h.Remark,
	}, nil
}

func (c *Convertor) ConvertPARowPoToDo(row *ent.PARow) *payment_center.PARow {
	return &payment_center.PARow{
		HeadCode:       row.HeadCode,
		RowCode:        row.RowCode,
		GrnCount:       row.GrnCount,
		GrnAmount:      row.GrnAmount,
		PayAmount:      row.PayAmount,
		DocDescription: row.Description,
	}
}
