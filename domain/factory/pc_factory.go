package factory

import (
	"context"

	"github.com/google/uuid"
	"github.com/spf13/cast"

	"purchase/domain/entity/payment_center"
	"purchase/domain/sal"
)

func NewPCFactory(mdm sal.MDMService) *PCFactory {
	return &PCFactory{
		mdm: mdm,
	}
}

type PCFactory struct {
	mdm sal.MDMService
}

func (f *PCFactory) BuildPA(ctx context.Context, h *payment_center.PAHead) (*payment_center.PAHead, error) {
	code := h.Code
	if code == "" {
		code = f.generateCode()
	}
	applicant, err := f.mdm.GetUser(ctx, h.Applicant.Account)
	if err != nil {
		return nil, err
	}

	dept, err := f.mdm.GetDepartment(ctx, h.Department.Code)
	if err != nil {
		return nil, err
	}

	rows := h.Rows
	for i, row := range rows {
		row.HeadCode = code
		row.RowCode = code + cast.ToString(i+1)
	}
	head := &payment_center.PAHead{
		Code:       code,
		State:      h.State,
		PayAmount:  h.PayAmount,
		Applicant:  applicant,
		Department: dept,
		Rows:       rows,
	}
	return head, nil
}

func (f *PCFactory) UpdatePA(ctx context.Context, pa, update *payment_center.PAHead) error {

	if pa.Applicant.Account != update.Applicant.Account {
		applicant, err := f.mdm.GetUser(ctx, update.Applicant.Account)
		if err != nil {
			return err
		}
		pa.Applicant = applicant
	}

	if pa.Department.Code != update.Department.Code {
		dept, err := f.mdm.GetDepartment(ctx, update.Department.Code)
		if err != nil {
			return err
		}
		pa.Department = dept
	}

	pa.State = update.State
	pa.PayAmount = update.PayAmount
	pa.Rows = update.Rows
	return nil
}

func (f *PCFactory) generateCode() string {
	return uuid.New().String()
}
