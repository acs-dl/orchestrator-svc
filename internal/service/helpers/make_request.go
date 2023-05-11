package helpers

import (
	"encoding/json"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func MakeGetRoleRequest(params data.RequestParams) (*resources.RoleResponse, error) {
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

	var response resources.RoleResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal body")
	}

	return &response, nil
}

func MakeCheckSubmoduleRequest(params data.RequestParams) (*resources.LinkResponse, error) {
	res, err := MakeHttpRequest(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make http request")
	}

	res, err = HandleHttpResponseStatusCode(res, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check response status code")
	}
	if res == nil {
		return nil, errors.Wrap(err, "something wring with response")
	}

	var response resources.LinkResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal body")
	}

	return &response, nil
}

func MakeNoResponseRequest(params data.RequestParams) error {
	res, err := MakeHttpRequest(params)
	if err != nil {
		return errors.Wrap(err, "failed to make http request")
	}

	res, err = HandleHttpResponseStatusCode(res, params)
	if err != nil {
		return errors.Wrap(err, "failed to check response status code")
	}
	if res == nil {
		return errors.New("something wrong with response")
	}

	return nil
}

func MakeGetUserRequest(params data.RequestParams, counter int64) (*resources.User, error) {
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

func MakeGetEstimatedTimeRequest(params data.RequestParams) (*resources.EstimatedTimeResponse, error) {
	res, err := MakeHttpRequest(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make http request")
	}

	res, err = HandleHttpResponseStatusCode(res, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check response status code")
	}
	if res == nil {
		return nil, errors.Wrap(err, "estimated time wasn't found")
	}

	var response resources.EstimatedTimeResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal body")
	}

	return &response, nil
}

func MakeGetRolesRequest(params data.RequestParams) (*data.ModuleRolesResponse, error) {
	res, err := MakeHttpRequest(params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make http request")
	}

	res, err = HandleHttpResponseStatusCode(res, params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to check response status code")
	}
	if res == nil {
		return nil, errors.New("something wrong with response")
	}

	var response data.ModuleRolesResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal body")
	}

	return &response, nil
}
