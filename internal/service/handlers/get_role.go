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

	roleResponse, err := helpers.MakeGetRoleRequest(data.RequestParams{
		Method: http.MethodGet,
		Link:   module.Link + "/role",
		Body:   nil,
		Header: map[string]string{
			"Content-Type": "application/json",
		},
		Query: map[string]string{
			"filter[accessLevel]": *request.AccessLevel,
		},
		Timeout: 30 * time.Second,
	})
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
