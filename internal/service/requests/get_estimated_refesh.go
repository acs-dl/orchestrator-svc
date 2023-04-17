package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetEstimatedRefreshRequest struct {
	ModuleName *string `filter:"moduleName"`
	Submodule  *string `filter:"submodule"`
}

func NewGetEstimatedRefreshRequest(r *http.Request) (GetEstimatedRefreshRequest, error) {
	var request GetEstimatedRefreshRequest

	err := urlval.Decode(r.URL.Query(), &request)

	return request, err
}
