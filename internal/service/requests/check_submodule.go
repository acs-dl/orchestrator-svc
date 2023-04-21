package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type CheckSubmoduleRequest struct {
	ModuleName *string `filter:"moduleName"`
	Submodule  *string `filter:"submodule"`
}

func NewCheckSubmoduleRequest(r *http.Request) (CheckSubmoduleRequest, error) {
	var request CheckSubmoduleRequest

	err := urlval.Decode(r.URL.Query(), &request)

	return request, err
}
