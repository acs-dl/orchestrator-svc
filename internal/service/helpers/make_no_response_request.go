package helpers

import (
	"bytes"
	"net/http"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func MakeNoResponseRequest(params data.RequestParams) error {
	req, err := http.NewRequest(params.Method, params.Link, bytes.NewReader(params.Body))
	if err != nil {
		return errors.Wrap(err, "couldn't create request")
	}

	req.Header.Set("Content-Type", "application/json")

	if params.AuthHeader != nil {
		req.Header.Set("Authorization", *params.AuthHeader)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error making http request")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.Errorf("error in response, status %s", res.Status)
	}

	return nil
}
