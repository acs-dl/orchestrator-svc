package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func MakeGetUserRequest(moduleLink, userId string, counter int64) (*resources.User, error) {
	link := fmt.Sprintf(moduleLink+"/users/%s", userId)
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create request")
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error making http request")
	}

	if res.StatusCode == 404 {
		return nil, nil
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("error in response, status %s", res.Status))
	}

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

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, " failed to unmarshal body")
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
