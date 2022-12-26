package requests

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type DeleteModuleRequest struct {
	ModuleName string `filter:"module_name"`
}

func NewDeleteModuleRequest(r *http.Request) (DeleteModuleRequest, error) {
	var request DeleteModuleRequest

	err := urlval.Decode(r.URL.Query(), &request)
	return request, errors.Wrap(err, "failed to parse request")
}
