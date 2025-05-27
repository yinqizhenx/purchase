package dtx

type StepStatus string

func (s StepStatus) String() string {
	return string(s)
}

const (
	StepPending           StepStatus = "pending"
	StepInAction          StepStatus = "in_action"
	StepActionFail        StepStatus = "action_fail"
	StepActionSuccess     StepStatus = "action_success"
	StepInCompensate      StepStatus = "in_compensate"
	StepCompensateSuccess StepStatus = "compensate_success"
	StepCompensateFail    StepStatus = "compensate_fail"
	// StepSuccess           StepStatus = 7
	// StepFailed            StepStatus = 8
)

type SagaState string

const (
	SagaStateExecuting SagaState = "executing"
	SagaStateSuccess   SagaState = "success"
	SagaStateFail      SagaState = "fail"
)
