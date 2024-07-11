package factory

import (
	"context"

	"github.com/google/uuid"
	"github.com/spf13/cast"

	"purchase/domain/entity/payment_center"
	"purchase/domain/sal"
	"purchase/domain/vo"
	pb "purchase/idl/payment_center"
)

func NewPCFactory(mdm sal.MDMService) *PCFactory {
	return &PCFactory{
		mdm: mdm,
	}
}

type PCFactory struct {
	mdm sal.MDMService
}

func (f *PCFactory) BuildPA(ctx context.Context, dto *pb.AddPAReq) (*payment_center.PAHead, error) {
	code := dto.Code
	if code == "" {
		code = f.generateCode()
	}

	state := vo.DocStateDraft
	if dto.IsSubmit {
		state = vo.DocStateSubmitted
	}

	applicant, err := f.mdm.GetUser(ctx, dto.Applicant)
	if err != nil {
		return nil, err
	}

	dept, err := f.mdm.GetDepartment(ctx, dto.ApplyDepartment)
	if err != nil {
		return nil, err
	}

	// sup, err := f.mdm.GetSupplier(ctx, dto.SupplierCode)
	// if err != nil {
	// 	return nil, err
	// }

	head := &payment_center.PAHead{
		Code:       code,
		State:      state,
		PayAmount:  dto.PayAmount,
		Applicant:  applicant,
		Department: dept,
		// Supplier:   sup,
	}

	for i, row := range dto.Rows {
		r := &payment_center.PARow{
			HeadCode:       dto.Code,
			RowCode:        code + cast.ToString(i+1),
			GrnCount:       row.GrnCount,
			GrnAmount:      row.GrnAmount,
			PayAmount:      row.PayAmount,
			DocDescription: row.DocDescription,
		}
		head.Rows = append(head.Rows, r)
	}

	return head, f.Validate(head)
}

func (f *PCFactory) UpdateBuildPA(ctx context.Context, pa *payment_center.PAHead, update *pb.UpdatePAReq) error {
	state := vo.DocStateDraft
	if update.IsSubmit {
		state = vo.DocStateSubmitted
	}

	applicant, err := f.mdm.GetUser(ctx, update.Applicant)
	if err != nil {
		return err
	}
	pa.Applicant = applicant

	dept, err := f.mdm.GetDepartment(ctx, update.ApplyDepartment)
	if err != nil {
		return err
	}

	sup, err := f.mdm.GetSupplier(ctx, update.SupplierCode)
	if err != nil {
		return err
	}

	pa.Department = dept
	pa.State = state
	pa.PayAmount = update.PayAmount
	pa.Applicant = applicant
	pa.Department = dept
	pa.Supplier = sup
	pa.Rows = nil // 清空

	for _, pbRow := range update.Rows {
		pa.Rows = append(pa.Rows, &payment_center.PARow{
			HeadCode:       pbRow.Code,
			RowCode:        pbRow.RowCode,
			GrnCount:       pbRow.GrnCount,
			GrnAmount:      pbRow.GrnAmount,
			PayAmount:      pbRow.PayAmount,
			DocDescription: pbRow.DocDescription,
		})
	}
	return f.Validate(pa)
}

func (f *PCFactory) generateCode() string {
	return uuid.New().String()
}

func (f *PCFactory) Validate(h *payment_center.PAHead) error {
	return nil
}
