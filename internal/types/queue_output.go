package types

import "gitlab.com/distributed_lab/acs/orchestrator/internal/data"

type ModuleResult string

const (
	ModuleResultSuccess ModuleResult = "success"
	ModuleResultInvited ModuleResult = "invited"
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
	switch mr {
	case ModuleResultSuccess:
		return data.FINISHED
	case ModuleResultInvited:
		return data.INVITED
	default:
		return data.FAILED
	}
}
