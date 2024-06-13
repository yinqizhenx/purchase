package assembler

import (
	"purchase/domain/entity/payment_center"
	"purchase/domain/entity/user"
	pb "purchase/idl/payment_center"
)

type PaymentAssembler struct{}

func NewPaymentAssembler() *PaymentAssembler {
	return &PaymentAssembler{}
}

func (asb PaymentAssembler) PAHeadDoToDto() (*pb.AddPAReq, error) {
	return nil, nil
}

func (asb PaymentAssembler) PAHeadDtoToDo(dto *pb.AddPAReq) *payment_center.PAHead {
	h := &payment_center.PAHead{
		PayAmount: dto.PayAmount,
		Code:      dto.Code,
		State:     dto.State,
		Applicant: user.User{
			Account: dto.Applicant,
		},
	}
	return h
}
