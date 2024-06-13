package vo

type ContractOnceRowStage string

const (
	ContractRowStageTotal ContractOnceRowStage = ""     // 所有阶段的汇总值
	ContractRowStagePre   ContractOnceRowStage = "pre"  // 预付款阶段
	ContractRowStagePost  ContractOnceRowStage = "post" // 后付款阶段
)
