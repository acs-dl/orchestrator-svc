package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"net/http"
)

type CreateRequestRequest struct {
	Data resources.Request `json:"data"`
}

func NewCreateRequestRequest(r *http.Request) (CreateRequestRequest, error) {
	var request CreateRequestRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}

	return request, request.validate()
}

func (r *CreateRequestRequest) validate() error {
	return validation.Errors{
		"module":    validation.Validate(&r.Data.Attributes.Module, validation.Required),
		"payload":   validation.Validate(&r.Data.Attributes.Payload, validation.Required),
		"from_user": validation.Validate(&r.Data.Attributes.FromUser, validation.Required),
		"to_user":   validation.Validate(&r.Data.Attributes.ToUser, validation.Required),
	}.Filter()
}
