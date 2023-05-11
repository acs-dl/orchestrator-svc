package helpers

import (
	"encoding/json"

	"github.com/acs-dl/orchestrator-svc/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func CreateJsonSubmodulesBody(submodules []string) ([]byte, error) {
	type req struct {
		Data resources.Submodule `json:"data"`
	}

	body := resources.Submodule{
		Key: resources.Key{
			ID:   string(resources.LINKS),
			Type: resources.LINKS,
		},
		Attributes: resources.SubmoduleAttributes{
			Links: submodules,
		},
	}

	newReq := req{Data: body}

	bodyBytes, err := json.Marshal(newReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal body")
	}

	return bodyBytes, nil
}
