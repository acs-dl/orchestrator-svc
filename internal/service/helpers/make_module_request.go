package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func MakeRequest(moduleLink, userId string, counter int64) (*resources.User, error) {
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

	returned := struct {
		Data struct {
			Attributes struct {
				Module   string  `json:"module"`
				UserId   int64   `json:"user_id"`
				Username *string `json:"username"`
				Phone    *string `json:"phone"`
			} `json:"attributes"`
		} `json:"data"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&returned); err != nil {
		return nil, errors.Wrap(err, " failed to unmarshal body")
	}

	return &resources.User{
		Key: resources.NewKeyInt64(counter, resources.USERS),
		Attributes: resources.UserAttributes{
			Module:   returned.Data.Attributes.Module,
			UserId:   returned.Data.Attributes.UserId,
			Username: returned.Data.Attributes.Username,
			Phone:    returned.Data.Attributes.Phone,
		},
	}, nil
}
