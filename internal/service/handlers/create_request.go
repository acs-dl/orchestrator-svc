package handlers

import (
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strconv"
)

func CreateRequest(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateRequestRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse create request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	toUserId, err := strconv.Atoi(request.Request.Relationships.User.Data.ID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse user id")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	module := data.Module{Name: request.Request.Attributes.Module}
	requestData := data.Request{
		ID: uuid.New().String(),
		// TODO: add from user id
		FromUserID: 12345,
		ToUserID:   int64(toUserId),
		Payload:    request.Request.Attributes.Payload,
		ModuleName: module.Name,
		Status:     data.CREATED,
	}

	err = helpers.RequestsQ(r).Insert(requestData)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to save new request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	helpers.Log(r).Infof("successfully created request with id `%s`", requestData.ID)
	ape.Render(w, newCreateRequestResponse(requestData))
}

func newCreateRequestResponse(request data.Request) resources.Request {
	return resources.Request{
		Key: resources.Key{
			ID:   request.ID,
			Type: resources.REQUESTS,
		},
		Attributes: resources.RequestAttributes{
			Module:  request.ModuleName,
			Payload: request.Payload,
			Status:  string(request.Status),
		},
		Relationships: resources.RequestRelationships{},
	}
}
