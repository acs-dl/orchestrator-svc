package handlers

import (
	"encoding/json"
	"fmt"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func GetUserById(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetUserByIdRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse get user by id request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	modules, err := helpers.ModulesQ(r).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to select modules")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if len(modules) == 0 {
		helpers.Log(r).WithError(err).Error("no modules were found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var result []resources.User
	for _, module := range modules {
		returned, err := makeRequest(module.Name, request.Id)
		if err != nil {
			helpers.Log(r).WithError(err).Errorf("failed to get user with id `%s` from module `%s`", request.Id, module.Name)
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if returned == nil {
			continue
		}
		result = append(result, *returned)
	}

	ape.Render(w, newGetUserResponse(result))
}

func makeRequest(moduleName, userId string) (*resources.User, error) {
	link := fmt.Sprintf("http://localhost:9000/integrations/%s/users/%s", moduleName, userId)
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
				Module   string `json:"module"`
				UserId   int64  `json:"user_id"`
				Username string `json:"username"`
			} `json:"attributes"`
		} `json:"data"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&returned); err != nil {
		return nil, errors.Wrap(err, " failed to unmarshal body")
	}

	return &resources.User{
		Key: resources.NewKeyInt64(returned.Data.Attributes.UserId, resources.USERS),
		Attributes: resources.UserAttributes{
			Module:   returned.Data.Attributes.Module,
			UserId:   returned.Data.Attributes.UserId,
			Username: returned.Data.Attributes.Username,
		},
	}, nil
}

func newGetUserResponse(user []resources.User) resources.UserListResponse {
	return resources.UserListResponse{
		Data: user,
	}
}
