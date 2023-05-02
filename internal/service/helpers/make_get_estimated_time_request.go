package helpers

import (
	"encoding/json"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

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
