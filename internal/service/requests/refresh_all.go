package requests

import (
	"encoding/json"
	"net/http"

	"github.com/acs-dl/orchestrator-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RefreshAllRequest struct {
	Data resources.FromToUser `json:"data"`
}

func NewRefreshAllRequest(r *http.Request) (RefreshAllRequest, error) {
	var request RefreshAllRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}

	return request, request.validate()
}

func (r *RefreshAllRequest) validate() error {
	return validation.Errors{
		"from_user": validation.Validate(&r.Data.Attributes.FromUser, validation.Required),
		"to_user":   validation.Validate(&r.Data.Attributes.ToUser, validation.Required),
	}.Filter()
}
