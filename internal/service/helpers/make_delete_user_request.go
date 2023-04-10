package helpers

import (
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func MakeDeleteUserRequest(moduleLink, userId, authHeader string) error {
	link := fmt.Sprintf(moduleLink+"/users/%s", userId)
	req, err := http.NewRequest(http.MethodDelete, link, nil)
	if err != nil {
		return errors.Wrap(err, "couldn't create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error making http request")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("error in response, status %s", res.Status))
	}

	return nil
}
