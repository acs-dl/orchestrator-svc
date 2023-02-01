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

	toUserId, err := strconv.ParseInt(request.Data.Relationships.User.Data.ID, 10, 64)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse user id")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	requestData := data.Request{
		ID: uuid.New().String(),
		// TODO: add from user id
		FromUserID: 12345,
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

	helpers.Log(r).Infof("successfully created request with id `%s`", requestData.ID)
	ape.Render(w, newRequestResponse(requestData))
}

func newRequest(request data.Request) resources.Request {
	key := resources.NewKeyInt64(request.ToUserID, resources.USERS)
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
		Relationships: resources.RequestRelationships{
			User: resources.Relation{
				Data: &key,
			},
		},
	}
}

func newRequestResponse(request data.Request) resources.RequestResponse {
	return resources.RequestResponse{
		Data: newRequest(request),
	}
}
