package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type CheckSubmoduleRequest struct {
	ModuleName *string `filter:"module_name"`
	Submodule  *string `filter:"submodule"`
}

func NewCheckSubmoduleRequest(r *http.Request) (CheckSubmoduleRequest, error) {
	var request CheckSubmoduleRequest

	err := urlval.Decode(r.URL.Query(), &request)

	return request, err
}
