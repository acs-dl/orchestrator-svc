package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/acs-dl/orchestrator-svc/internal/data"
	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"github.com/acs-dl/orchestrator-svc/internal/service/requests"
	"github.com/acs-dl/orchestrator-svc/resources"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewDeleteUserByIdRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse get user by id request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	userId, err := strconv.ParseInt(request.Data.Attributes.ToUser, 10, 64)
	if err != nil {
		helpers.Log(r).WithError(err).Errorf("failed to parse to user id `%s`", request.Data.Attributes.ToUser)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	fromUserId, err := strconv.ParseInt(request.Data.Attributes.ToUser, 10, 64)
	if err != nil {
		helpers.Log(r).WithError(err).Errorf("failed to parse from user id `%s`", request.Data.Attributes.FromUser)
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

	var userinfoModules = make([]resources.User, 0)
	for i, module := range modules {
		response, err := helpers.MakeGetUserRequest(data.RequestParams{
			Method: http.MethodGet,
			Link:   fmt.Sprintf(module.Link+"/users/%s", userId),
			Header: map[string]string{
				"Content-Type": "application/json",
			},
			Body:    nil,
			Query:   nil,
			Timeout: 30 * time.Second,
		}, int64(i))
		if err != nil {
			helpers.Log(r).WithError(err).Errorf("failed to get user with id `%s` from module `%s`", request.Data.Attributes.ToUser, module.Name)
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if response == nil {
			continue
		}
		userinfoModules = append(userinfoModules, *response)
	}

	var requestToCheck map[string]bool
	for _, userinfoModule := range userinfoModules {
		module, err := helpers.ModulesQ(r).FilterByNames(userinfoModule.Attributes.Module).Get()
		if err != nil {
			helpers.Log(r).WithError(err).Errorf("failed to get module with name `%s`", userinfoModule.Attributes.Module)
			ape.RenderErr(w, problems.InternalError())
			return
		}

		if module == nil {
			helpers.Log(r).WithError(err).Errorf("no module with name `%s`", userinfoModule.Attributes.Module)
			ape.RenderErr(w, problems.NotFound())
			return
		}

		deleteUserJson, err := json.Marshal(data.DeleteUserPayload{
			Action:   data.DeleteUserAction,
			Username: userinfoModule.Attributes.Username,
			Phone:    userinfoModule.Attributes.Phone,
		})
		if err != nil {
			helpers.Log(r).WithError(err).Errorf("failed to marshal delete user payload")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		requestData := data.Request{
			ID:         uuid.New().String(),
			FromUserID: fromUserId,
			ToUserID:   userId,
			Payload:    deleteUserJson,
			ModuleName: module.Name,
			Status:     data.CREATED,
		}

		err = helpers.RequestsQ(r).Insert(requestData)
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to save new request")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		requestToCheck[requestData.ID] = false
		helpers.Log(r).Infof("successfully created request with id `%s`", requestData.ID)
	}

	marshalledRequests, err := json.Marshal(requestToCheck)
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

	helpers.Log(r).Infof("successfully created requests to delete user with id `%s` from modules", request.Data.Attributes.ToUser)
	ape.Render(w, http.StatusAccepted)
}
