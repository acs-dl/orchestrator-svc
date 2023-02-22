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

func GetRole(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRoleRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if request.ModuleName == nil || request.AccessLevel == nil {
		helpers.Log(r).WithError(err).Info("some filters are empty")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	module, err := helpers.ModulesQ(r).FilterByNames(*request.ModuleName).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Infof("failed to get module `%s`", *request.ModuleName)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if module == nil {
		helpers.Log(r).Errorf("no such module `%s`", *request.ModuleName)
		ape.Render(w, problems.NotFound())
		return
	}

	roleResponse, err := makeGetRoleRequest(module.Link, *request.AccessLevel)
	if err != nil {
		helpers.Log(r).WithError(err).Infof("failed to get role `%s` from module `%s`", *request.AccessLevel, *request.ModuleName)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if roleResponse == nil {
		helpers.Log(r).Errorf("no such role `%s`", *request.AccessLevel)
		ape.Render(w, problems.NotFound())
		return
	}

	ape.Render(w, roleResponse)
}

func makeGetRoleRequest(moduleLink, accessLevel string) (*resources.RoleResponse, error) {
	link := fmt.Sprintf(moduleLink + "/role")
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create request")
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("filter[accessLevel]", accessLevel)
	req.URL.RawQuery = q.Encode()

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

	var returned resources.RoleResponse

	if err := json.NewDecoder(res.Body).Decode(&returned); err != nil {
		return nil, errors.Wrap(err, " failed to unmarshal body")
	}

	return &returned, nil
}
