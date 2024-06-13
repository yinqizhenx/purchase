package assembler

import (
	"purchase/domain/entity/su"
	pb "purchase/idl/payment_center"
)

type SuAssembler struct{}

func NewSuAssembler() *SuAssembler {
	return &SuAssembler{}
}

func (asb *SuAssembler) SuToDispatchPB(su *su.SU) *pb.AddPARes {
	return nil
}
