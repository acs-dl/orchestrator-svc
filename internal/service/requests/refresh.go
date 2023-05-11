package requests

import (
	"encoding/json"
	"net/http"

	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type RefreshRequest struct {
	Data resources.Refresh `json:"data"`
}

func NewRefreshRequest(r *http.Request) (RefreshRequest, error) {
	var request RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, nil
}
