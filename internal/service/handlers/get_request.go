package handlers

import (
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetRequest(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRequestRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse get request request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	req, err := helpers.RequestsQ(r).FilterByIDs(request.Uuid).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get request")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if req == nil {
		helpers.Log(r).Errorf("no request with such id")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, newGetRequestResponse(*req))
}

func newGetRequestResponse(request data.Request) resources.Request {
	return resources.Request{
		Key: resources.Key{
			ID:   request.ID,
			Type: resources.REQUESTS,
		},
		Attributes: resources.RequestAttributes{
			Module:  request.ModuleName,
			Payload: request.Payload,
			Status:  string(request.Status),
			Error:   request.Error,
		},
		Relationships: resources.RequestRelationships{},
	}
}
