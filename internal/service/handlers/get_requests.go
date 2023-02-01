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

func GetRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRequestsRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse get request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	requestsQ := helpers.RequestsQ(r)
	if request.FromUserId != nil {
		requestsQ = requestsQ.FilterByFromIds(*request.FromUserId)
	}
	if request.ToUserId != nil {
		requestsQ = requestsQ.FilterByToIds(*request.ToUserId)
	}

	dbRequests, err := requestsQ.Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get requests")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newRequestArrayResponse(dbRequests))
}

func newRequestArrayResponse(requests []data.Request) resources.RequestListResponse {
	var result []resources.Request

	for _, request := range requests {
		result = append(result, newRequest(request))
	}

	return resources.RequestListResponse{
		Data: result,
	}
}
