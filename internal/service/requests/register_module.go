package requests

import (
	"encoding/json"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"net/http"
)

type RegisterModuleRequest struct {
	Data resources.Module `json:"data"`
}

func NewRegisterModuleRequest(r *http.Request) (RegisterModuleRequest, error) {
	var request RegisterModuleRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}

	return request, nil
}
