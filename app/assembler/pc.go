package assembler

import (
	"purchase/domain/company"
	"purchase/domain/entity/department"
	"purchase/domain/entity/payment_center"
	"purchase/domain/entity/supplier"
	"purchase/domain/entity/user"
	"purchase/domain/vo"
	pb "purchase/idl/payment_center"
)

func (a *Assembler) PAHeadDoToDto() (*pb.AddOrUpdatePAReq, error) {
	return nil, nil
}

func (a *Assembler) PAHeadDtoToDo(dto *pb.AddOrUpdatePAReq) *payment_center.PAHead {
	state := vo.DocStateDraft
	if dto.IsSubmit {
		state = vo.DocStateSubmitted
	}
	h := &payment_center.PAHead{
		PayAmount: dto.PayAmount,
		Code:      dto.Code,
		State:     state,
		Applicant: user.User{
			Account: dto.Applicant,
		},
		Department: department.Department{
			Code: dto.ApplyDepartment,
		},
		Currency:   dto.Currency,
		IsAdv:      dto.IsAdv,
		HasInvoice: dto.HasInvoice,
		Company: company.Company{
			Code: dto.Code,
		},
		Supplier: supplier.Supplier{},
		Remark:   dto.Remark,
	}
	for _, row := range dto.Rows {
		r := &payment_center.PARow{
			HeadCode: dto.Code,
			// RowCode        string
			GrnCount:       row.GrnCount,
			GrnAmount:      row.GrnAmount,
			PayAmount:      row.PayAmount,
			DocDescription: row.DocDescription,
		}
		h.Rows = append(h.Rows, r)
	}
	return h
}
