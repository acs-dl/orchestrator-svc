package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func CheckSubmodule(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCheckSubmoduleRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if request.ModuleName == nil || request.Submodule == nil {
		helpers.Log(r).WithError(err).Info("some filters are empty")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	module, err := helpers.ModulesQ(r).FilterByNames(*request.ModuleName).FilterByIsModule(true).Get()
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

	submoduleResponse, err := makeCheckSubmoduleRequest(module.Link, *request.Submodule, r.Header.Get("Authorization"))
	if err != nil {
		helpers.Log(r).WithError(err).Infof("failed to check submodule `%s` from module `%s`", *request.Submodule, *request.ModuleName)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if submoduleResponse == nil {
		helpers.Log(r).Errorf("wrong submodule respone `%s`", *request.Submodule)
		ape.Render(w, problems.NotFound())
		return
	}

	ape.Render(w, submoduleResponse)
}

func makeCheckSubmoduleRequest(moduleLink, submodule, authHeader string) (*resources.LinkResponse, error) {
	link := fmt.Sprintf(moduleLink + "/submodule")
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	q := req.URL.Query()
	q.Add("filter[link]", submodule)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error making http request")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("error in response, status %s", res.Status))
	}

	var response resources.LinkResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, " failed to unmarshal body")
	}

	return &response, nil
}
