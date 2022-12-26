package requests

import (
	"encoding/json"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"net/http"
)

type CreateRequestRequest struct {
	Request resources.Request `json:"request"`
}

func NewCreateRequestRequest(r *http.Request) (CreateRequestRequest, error) {
	var request CreateRequestRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}

	return request, nil
}
