package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetUserByIdRequest struct {
	Id string
}

func NewGetUserByIdRequest(r *http.Request) (GetUserByIdRequest, error) {
	var request GetUserByIdRequest

	request.Id = chi.URLParam(r, "id")

	return request, request.validate()
}

func (r *GetUserByIdRequest) validate() error {
	return validation.Errors{
		"id": validation.Validate(&r.Id, validation.Required),
	}.Filter()
}
