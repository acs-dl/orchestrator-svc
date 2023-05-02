package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func MakeHttpRequest(params data.RequestParams) (*data.ResponseParams, error) {
	req, err := http.NewRequest(params.Method, params.Link, bytes.NewReader(params.Body))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create request")
	}

	ctx, cancel := context.WithTimeout(context.Background(), params.Timeout)
	defer cancel()
	req = req.WithContext(ctx)

	if params.Header != nil {
		for key, value := range params.Header {
			req.Header.Set(key, value)
		}
	}

	if params.Query != nil {
		q := req.URL.Query()
		for key, value := range params.Query {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error making http request")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading response body")
	}

	return &data.ResponseParams{
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     response.Header,
		StatusCode: response.StatusCode,
	}, nil
}

func HandleHttpResponseStatusCode(response *data.ResponseParams, params data.RequestParams) (*data.ResponseParams, error) {
	switch status := response.StatusCode; {
	case status >= http.StatusOK && status < http.StatusMultipleChoices:
		return response, nil
	case status == http.StatusNotFound:
		return nil, nil
	case status < http.StatusOK || status >= http.StatusMultipleChoices:
		return nil, errors.New(fmt.Sprintf("error in response `%s`", http.StatusText(response.StatusCode)))
	default:
		return nil, errors.New(fmt.Sprintf("unexpected response status `%s`", http.StatusText(response.StatusCode)))
	}
}
