package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type GetRequestRequest struct {
	Uuid string `filter:"id"`
}

func NewGetRequestRequest(r *http.Request) (GetRequestRequest, error) {
	var request GetRequestRequest

	request.Uuid = chi.URLParam(r, "id")

	return request, request.validate()
}

func (r *GetRequestRequest) validate() error {
	return validation.Errors{
		"uuid": validation.Validate(&r.Uuid, validation.Required),
	}.Filter()
}
