package requests

import (
	"encoding/json"
	"net/http"

	"github.com/acs-dl/orchestrator-svc/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DeleteUserByIdRequest struct {
	Data resources.Request `json:"data"`
}

func NewDeleteUserByIdRequest(r *http.Request) (*DeleteUserByIdRequest, error) {
	var request DeleteUserByIdRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return &request, request.validate()
}

func (r *DeleteUserByIdRequest) validate() error {
	return validation.Errors{
		"to_user_id":   validation.Validate(&r.Data.Attributes.ToUser, validation.Required),
		"from_user_id": validation.Validate(&r.Data.Attributes.FromUser, validation.Required),
	}.Filter()
}
