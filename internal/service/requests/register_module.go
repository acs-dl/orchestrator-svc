package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
)

type RegisterModuleRequest struct {
	Data resources.Module `json:"data"`
}

func NewRegisterModuleRequest(r *http.Request) (RegisterModuleRequest, error) {
	var request RegisterModuleRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}

	return request, request.validate()
}

func (r *RegisterModuleRequest) validate() error {
	return validation.Errors{
		"name":      validation.Validate(&r.Data.Attributes.Name, validation.Required),
		"topic":     validation.Validate(&r.Data.Attributes.Topic, validation.Required),
		"link":      validation.Validate(&r.Data.Attributes.Link, validation.Required),
		"title":     validation.Validate(&r.Data.Attributes.Title, validation.Required),
		"is_module": validation.Validate(&r.Data.Attributes.IsModule, validation.In(true, false)),
	}.Filter()
}
