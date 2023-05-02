package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func MakeGetUserRequest(moduleLink, userId string, counter int64) (*resources.User, error) {
	params := data.RequestParams{
		Method: http.MethodGet,
		Link:   fmt.Sprintf(moduleLink+"/users/%s", userId),
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		Body:    nil,
		Query:   nil,
		Timeout: 30 * time.Second,
	}

	res, err := MakeHttpRequest(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make http request")
	}

	res, err = HandleHttpResponseStatusCode(res, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check response status code")
	}
	if res == nil {
		return nil, nil
	}

	return populateGetUserResponse(res, counter)
}

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
