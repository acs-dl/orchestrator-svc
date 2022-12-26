package handlers

import (
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func RegisterModule(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewRegisterModuleRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse registration request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	module := data.Module{
		Name:     request.Data.Attributes.Name,
		Endpoint: &request.Data.Attributes.Endpoint,
	}

	err = helpers.ModulesQ(r).Insert(module)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to save new module")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusAccepted)
}
