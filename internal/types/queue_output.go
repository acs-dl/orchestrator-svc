package types

import "gitlab.com/distributed_lab/acs/orchestrator/internal/data"

type ModuleResult string

const (
	ModuleResultSuccess ModuleResult = "success"
	ModuleResultFailure ModuleResult = "failure"
)

type QueueOutput struct {
	ID        string       `json:"id"`
	Status    ModuleResult `json:"status"`
	Error     string       `json:"error"`
	RequestId *string      `json:"request_id,omitempty"`
	Action    *string      `json:"action,omitempty"`
}

func (mr ModuleResult) ToRequestStatus() data.RequestStatus {
	if mr == ModuleResultSuccess {
		return data.FINISHED
	}

	return data.FAILED
}
