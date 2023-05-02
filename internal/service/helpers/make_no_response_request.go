package helpers

import (
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

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
