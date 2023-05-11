package helpers

import (
	"encoding/json"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func populateGetUserResponse(res *data.ResponseParams, counter int64) (*resources.User, error) {
	response := struct {
		Data struct {
			Attributes struct {
				Module      string  `json:"module"`
				Submodule   string  `json:"submodule"`
				AccessLevel string  `json:"access_level"`
				UserId      int64   `json:"user_id"`
				Username    *string `json:"username"`
				Phone       *string `json:"phone"`
			} `json:"attributes"`
		} `json:"data"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal body")
	}

	return &resources.User{
		Key: resources.NewKeyInt64(counter, resources.USERS),
		Attributes: resources.UserAttributes{
			Module:      response.Data.Attributes.Module,
			UserId:      response.Data.Attributes.UserId,
			Username:    response.Data.Attributes.Username,
			Phone:       response.Data.Attributes.Phone,
			Submodule:   response.Data.Attributes.Submodule,
			AccessLevel: response.Data.Attributes.AccessLevel,
		},
	}, nil
}
