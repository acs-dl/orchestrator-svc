package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func RefreshAllModules(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewRefreshAllRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse create request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	fromUserId, err := strconv.ParseInt(request.Data.Attributes.FromUser, 10, 64)
	if err != nil {
		helpers.Log(r).WithError(err).Errorf("failed to parse from user id `%s`", request.Data.Attributes.FromUser)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	toUserId, err := strconv.ParseInt(request.Data.Attributes.ToUser, 10, 64)
	if err != nil {
		helpers.Log(r).WithError(err).Errorf("failed to parse to user id `%s`", request.Data.Attributes.ToUser)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	modules, err := helpers.ModulesQ(r).FilterByIsModule(true).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to select all modules")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	refreshModulePayload, err := helpers.CreateRefreshModulePayload()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create refresh module payload")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	for _, module := range modules {
		requestData := data.Request{
			ID:         uuid.New().String(),
			FromUserID: fromUserId,
			ToUserID:   toUserId,
			Payload:    refreshModulePayload,
			ModuleName: module.Name,
			Status:     data.CREATED,
		}

		err = helpers.RequestsQ(r).Insert(requestData)
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to save new request")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		helpers.Log(r).Infof("successfully created refresh request with id `%s`", requestData.ID)

		marshalledRequests, err := json.Marshal(map[string]bool{requestData.ID: false})
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to marshal requests")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		err = helpers.RequestTransactionsQ(r).Insert(data.RequestTransaction{
			ID:       uuid.New().String(),
			Action:   data.Single,
			Requests: marshalledRequests,
		})
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to save new request transaction")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	helpers.Log(r).Infof("successfully created requests to refresh all")
	w.WriteHeader(http.StatusAccepted)
	ape.Render(w, http.StatusAccepted)
}
