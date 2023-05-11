package handlers

import (
	"github.com/acs-dl/orchestrator-svc/internal/data"
	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"github.com/acs-dl/orchestrator-svc/internal/service/requests"
	"github.com/acs-dl/orchestrator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
	"strings"
)

func GetRequests(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRequestsRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse get request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	countRequestsQ := helpers.RequestsQ(r).Count()
	requestsQ := helpers.RequestsQ(r)
	if request.FromUserId != nil {
		requestsQ = requestsQ.FilterByFromIds(*request.FromUserId)
		countRequestsQ = countRequestsQ.FilterByFromIds(*request.FromUserId)
	}
	if request.ToUserId != nil {
		requestsQ = requestsQ.FilterByToIds(*request.ToUserId)
		countRequestsQ = countRequestsQ.FilterByToIds(*request.ToUserId)
	}
	if request.Status != nil {
		requestsQ = requestsQ.FilterByStatuses(data.RequestStatus(*request.Status))
		countRequestsQ = countRequestsQ.FilterByStatuses(data.RequestStatus(*request.Status))
	}
	if request.Action != nil {
		actions := strings.Split(*request.Action, " ")
		requestsQ = requestsQ.FilterNotByActions(actions...)
		countRequestsQ = countRequestsQ.FilterNotByActions(actions...)
	}

	dbRequests, err := requestsQ.Page(request.OffsetPageParams).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get requests")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	totalCount, err := countRequestsQ.GetTotalCount()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get total count requests")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := newRequestArrayResponse(dbRequests)
	response.Meta.TotalCount = totalCount
	response.Links = data.GetOffsetLinksForPGParams(r, request.OffsetPageParams)

	ape.Render(w, response)
}

func newRequestArrayResponse(requests []data.Request) RequestListResponse {
	result := make([]resources.Request, 0)

	for _, request := range requests {
		result = append(result, newRequest(request))
	}

	return RequestListResponse{
		Data: result,
	}
}

type RequestListResponse struct {
	Meta  Meta                `json:"meta"`
	Data  []resources.Request `json:"data"`
	Links *resources.Links    `json:"links"`
}

type Meta struct {
	TotalCount int64 `json:"total_count"`
}
