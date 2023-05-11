package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/acs-dl/orchestrator-svc/internal/data"
	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"github.com/acs-dl/orchestrator-svc/internal/service/requests"
	"github.com/acs-dl/orchestrator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetUserById(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetUserByIdRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse get user by id request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	modules, err := helpers.ModulesQ(r).FilterByIsModule(true).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to select modules")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if len(modules) == 0 {
		helpers.Log(r).WithError(err).Error("no modules were found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	var result = make([]resources.User, 0)
	for i, module := range modules {
		response, err := helpers.MakeGetUserRequest(data.RequestParams{
			Method: http.MethodGet,
			Link:   fmt.Sprintf(module.Link+"/users/%s", request.Id),
			Header: map[string]string{
				"Content-Type": "application/json",
			},
			Body:    nil,
			Query:   nil,
			Timeout: 30 * time.Second,
		}, int64(i))
		if err != nil {
			helpers.Log(r).WithError(err).Errorf("failed to get user with id `%s` from module `%s`", request.Id, module.Name)
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if response == nil {
			continue
		}
		result = append(result, *response)
	}

	ape.Render(w, newGetUserResponse(result))
}

func newGetUserResponse(user []resources.User) resources.UserListResponse {
	return resources.UserListResponse{
		Data: user,
	}
}
