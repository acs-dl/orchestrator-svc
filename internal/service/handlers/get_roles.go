package handlers

import (
	"net/http"
	"time"

	"github.com/acs-dl/orchestrator-svc/internal/data"
	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"github.com/acs-dl/orchestrator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	modules, err := helpers.ModulesQ(r).FilterByIsModule(true).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Errorf("failed to select modules")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := newModulesRolesResponse()

	if len(modules) == 0 {
		helpers.Log(r).Errorf("no modules was found")
		ape.Render(w, response)
		return
	}

	for _, module := range modules {
		moduleRoles, err := helpers.MakeGetRolesRequest(data.RequestParams{
			Method: http.MethodGet,
			Link:   module.Link + "/roles",
			Body:   nil,
			Header: map[string]string{
				"Content-Type": "application/json",
			},
			Query:   nil,
			Timeout: 30 * time.Second,
		})
		if err != nil {
			helpers.Log(r).WithError(err).Infof("failed to get roles from `%s`", module.Link+"/roles")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		response.Data.Attributes[module.Name] = moduleRoles.Data.Attributes
	}

	ape.Render(w, response)
}

func newModulesRolesResponse() data.ModulesRolesResponse {
	return data.ModulesRolesResponse{
		Data: data.ModulesRoles{
			Key: resources.Key{
				ID:   "0",
				Type: resources.MODULES,
			},
			Attributes: data.Modules{},
		},
	}
}
