package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
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

	requestToCheck := make([]string, 0)
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

		requestToCheck = append(requestToCheck, requestData.ID)
		helpers.Log(r).Infof("successfully created request with id `%s`", requestData.ID)
	}

	//requestToCheck - array with request ids for deleting user from modules

	err = waitForRequestsToHandleInModules(r, requestToCheck)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to wait request handling")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	err = checkIdentityRegisteredAndMakeDeleteUserRequest(helpers.ModulesQ(r), request.Data.Attributes.ToUser, r.Header.Get("Authorization"))
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to check identity and make delete request")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	helpers.Log(r).Infof("successfully created requests to delete user with id `%s` from modules", request.Data.Attributes.ToUser)
	ape.Render(w, http.StatusAccepted)
}

func waitForRequestsToHandleInModules(r *http.Request, requests []string) error {
	helpers.Log(r).Infof("started waiting to handle requests")

	for len(requests) != 0 {
		msgRequests, err := helpers.RequestsQ(r).FilterByIDs(requests...).Select()
		if err != nil {
			return errors.Wrap(err, "failed to select requests to check")
		}

		if len(msgRequests) == 0 {
			return errors.Errorf("no request returned")
		}

		for i, request := range msgRequests {
			if request.Status == data.FINISHED {
				requests = append(requests[:i], requests[i+1:]...)
				helpers.Log(r).Infof("request `%s` was handled `%d` more left", request.ID, len(requests))
				continue
			}
			if request.Status == data.FAILED {
				errMsg := ""
				if request.Error != nil {
					errMsg = *request.Error
				}
				return errors.Errorf("request `%s` returned error `%s`", request.ID, errMsg)
			}
		}

		time.Sleep(5 * time.Second)
	}

	helpers.Log(r).Infof("finished waiting to handle requests")
	return nil
}

func checkIdentityRegisteredAndMakeDeleteUserRequest(moduleQ data.ModuleQ, userId, authHeader string) error {
	module, err := moduleQ.FilterByNames("identity").Get()
	if err != nil {
		return errors.Wrap(err, "failed to get identity module")
	}

	if module == nil {
		return errors.Errorf("no module with name `identity`")
	}

	err = helpers.MakeNoResponseRequest(data.RequestParams{
		Method: http.MethodDelete,
		Link:   module.Link + "/users/" + userId,
		Header: map[string]string{
			"Authorization": authHeader,
			"Content-Type":  "application/json",
		},
		Body:    nil,
		Query:   nil,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		return errors.Wrap(err, "failed to make delete user request")
	}

	return nil
}
