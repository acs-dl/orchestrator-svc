package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/acs-dl/orchestrator-svc/internal/data"
	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"github.com/acs-dl/orchestrator-svc/internal/service/requests"
	"github.com/acs-dl/orchestrator-svc/resources"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateRequest(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateRequestRequest(r)
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

	module, err := helpers.ModulesQ(r).FilterByNames(request.Data.Attributes.Module).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Errorf("failed to get module with name `%s`", request.Data.Attributes.Module)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if module == nil {
		helpers.Log(r).WithError(err).Errorf("no module with name `%s`", request.Data.Attributes.Module)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	requestData := data.Request{
		ID:         uuid.New().String(),
		FromUserID: fromUserId,
		ToUserID:   toUserId,
		Payload:    request.Data.Attributes.Payload,
		ModuleName: request.Data.Attributes.Module,
		Status:     data.CREATED,
	}

	err = helpers.RequestsQ(r).Insert(requestData)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to save new request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

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

	helpers.Log(r).Infof("successfully created request with id `%s`", requestData.ID)
	w.WriteHeader(http.StatusAccepted)
	ape.Render(w, newRequestResponse(requestData))
}

func newRequest(request data.Request) resources.Request {
	return resources.Request{
		Key: resources.Key{
			ID:   request.ID,
			Type: resources.REQUESTS,
		},
		Attributes: resources.RequestAttributes{
			Module:    request.ModuleName,
			Payload:   request.Payload,
			CreatedAt: request.CreatedAt,
			Status:    string(request.Status),
			Error:     request.Error,
			FromUser:  fmt.Sprintf("%d", request.FromUserID),
			ToUser:    fmt.Sprintf("%d", request.ToUserID),
		},
		Relationships: resources.RequestRelationships{},
	}
}

func newRequestResponse(request data.Request) resources.RequestResponse {
	return resources.RequestResponse{
		Data: newRequest(request),
	}
}
