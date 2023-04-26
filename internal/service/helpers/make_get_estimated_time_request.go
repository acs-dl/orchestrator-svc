package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func MakeGetEstimatedTimeRequest(params data.RequestParams) (*resources.EstimatedTimeResponse, error) {
	req, err := http.NewRequest(params.Method, params.Link, bytes.NewReader(params.Body))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create request")
	}

	req.Header.Set("Content-Type", "application/json")

	if params.AuthHeader != nil {
		req.Header.Set("Authorization", *params.AuthHeader)
	}

	if params.Query != nil {
		q := req.URL.Query()
		for key, value := range params.Query {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error making http request")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.Errorf("error in response, status %s", res.Status)
	}

	var response resources.EstimatedTimeResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, " failed to unmarshal body")
	}

	return &response, nil
}
