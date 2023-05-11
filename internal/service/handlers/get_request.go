package handlers

import (
	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"github.com/acs-dl/orchestrator-svc/internal/service/requests"
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

	ape.Render(w, newRequestResponse(*req))
}
