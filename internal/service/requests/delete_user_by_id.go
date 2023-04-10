package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/urlval"
)

type DeleteUserByIdRequest struct {
	Id         string
	FromUserId *int64 `filter:"fromUserId"`
}

func NewDeleteUserByIdRequest(r *http.Request) (*DeleteUserByIdRequest, error) {
	var request DeleteUserByIdRequest

	request.Id = chi.URLParam(r, "id")

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return nil, err
	}

	return &request, request.validate()
}

func (r *DeleteUserByIdRequest) validate() error {
	return validation.Errors{
		"id":           validation.Validate(&r.Id, validation.Required),
		"from_user_id": validation.Validate(&r.FromUserId, validation.Required),
	}.Filter()
}
