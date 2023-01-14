package handlers

import (
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func UnregisterModule(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewDeleteModuleRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse unregister request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	err = helpers.ModulesQ(r).Delete(request.ModuleName)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete module")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	helpers.Log(r).Infof("successfully unregister module `%s`", request.ModuleName)
	w.WriteHeader(http.StatusAccepted)
}
