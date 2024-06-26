package assembler

import (
	"purchase/domain/entity/payment_center"
	"purchase/domain/entity/user"
	pb "purchase/idl/payment_center"
)

func (a *Assembler) PAHeadDoToDto() (*pb.AddPAReq, error) {
	return nil, nil
}

func (a *Assembler) PAHeadDtoToDo(dto *pb.AddPAReq) *payment_center.PAHead {
	h := &payment_center.PAHead{
		PayAmount: dto.PayAmount,
		Code:      dto.Code,
		Applicant: user.User{
			Account: dto.Applicant,
		},
	}
	return h
}
