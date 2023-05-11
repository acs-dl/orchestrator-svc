package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
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

	msgId := uuid.New().String()
	moduleNameJson, err := json.Marshal(data.PermissionsPayload{
		RequestId:  msgId,
		Action:     data.RemoveModuleAction,
		ModuleName: request.ModuleName,
	})
	if err != nil {
		helpers.Log(r).WithError(err).Errorf("failed to marshal module permission message")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = helpers.Sender(r).SendMessageToCustomChannel(data.AuthService, helpers.Sender(r).BuildPermissionsMessage(msgId, moduleNameJson))
	if err != nil {
		helpers.Log(r).WithError(err).Errorf("failed to send message to custom channel")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	helpers.Log(r).Infof("successfully unregister module `%s`", request.ModuleName)
	w.WriteHeader(http.StatusAccepted)
}
