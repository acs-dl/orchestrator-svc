package handlers

import (
	"net/http"
	"time"

	"github.com/acs-dl/orchestrator-svc/internal/data"
	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"github.com/acs-dl/orchestrator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
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

	submoduleResponse, err := helpers.MakeCheckSubmoduleRequest(data.RequestParams{
		Method: http.MethodGet,
		Link:   module.Link + "/submodule",
		Body:   nil,
		Header: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": r.Header.Get("Authorization"),
		},
		Query: map[string]string{
			"filter[link]": *request.Submodule,
		},
		Timeout: 30 * time.Second,
	})
	if err != nil {
		helpers.Log(r).WithError(err).Infof("failed to check submodule `%s` from module `%s`", *request.Submodule, *request.ModuleName)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if submoduleResponse == nil {
		helpers.Log(r).Errorf("wrong submodule response `%s`", *request.Submodule)
		ape.Render(w, problems.NotFound())
		return
	}

	ape.Render(w, submoduleResponse)
}
