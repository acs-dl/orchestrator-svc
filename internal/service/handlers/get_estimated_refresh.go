package handlers

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetEstimatedRefreshTime(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewRefreshRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if request.Data.Attributes.ModuleName == nil {
		estimatedTimeResponse, err := getEstimatedRefreshModules(helpers.ModulesQ(r), r.Header.Get("Authorization"))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to get estimated refresh modules")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, estimatedTimeResponse)
		return
	}

	if request.Data.Attributes.Submodule == nil {
		estimatedTimeResponse, err := getEstimatedRefreshModule(helpers.ModulesQ(r), *request.Data.Attributes.ModuleName, r.Header.Get("Authorization"))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to get estimated refresh module")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, estimatedTimeResponse)
		return
	}

	estimatedTimeResponse, err := getEstimatedRefreshModuleSubmodules(helpers.ModulesQ(r), *request.Data.Attributes.ModuleName, r.Header.Get("Authorization"), *request.Data.Attributes.Submodule)
	if err != nil {
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to get estimated refresh module submodule")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	ape.Render(w, estimatedTimeResponse)
}

func getEstimatedRefreshModuleSubmodules(modulesQ data.ModuleQ, moduleName, authHeader string, submodules []string) (*resources.EstimatedTimeResponse, error) {
	module, err := modulesQ.FilterByNames(moduleName).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get module")
	}

	if module == nil {
		return nil, errors.New("no such module")
	}

	body, err := helpers.CreateJsonSubmodulesBody(submodules)
	if err != nil {
		return nil, err
	}

	estimatedTime, err := helpers.MakeGetEstimatedTimeRequest(data.RequestParams{
		Method: http.MethodPost,
		Link:   module.Link + "/estimate_refresh/submodule",
		Header: map[string]string{
			"Authorization": authHeader,
		},
		Body:    body,
		Query:   nil,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make refresh request")
	}

	roundedTime, err := helpers.RoundDuration(estimatedTime.Data.Attributes.EstimatedTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to round duration")
	}

	estimatedTime.Data.Attributes.EstimatedTime = roundedTime.String()

	return estimatedTime, nil
}

func getEstimatedRefreshModule(modulesQ data.ModuleQ, moduleName, authHeader string) (*resources.EstimatedTimeResponse, error) {
	module, err := modulesQ.FilterByNames(moduleName).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get module")
	}

	if module == nil {
		return nil, errors.New("no such module")
	}

	estimatedTime, err := helpers.MakeGetEstimatedTimeRequest(data.RequestParams{
		Method: http.MethodPost,
		Link:   module.Link + "/estimate_refresh/module",
		Header: map[string]string{
			"Authorization": authHeader,
		},
		Body:    nil,
		Query:   nil,
		Timeout: 30 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make refresh request")
	}

	roundedTime, err := helpers.RoundDuration(estimatedTime.Data.Attributes.EstimatedTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to round duration")
	}

	estimatedTime.Data.Attributes.EstimatedTime = roundedTime.String()

	return estimatedTime, nil
}

func getEstimatedRefreshModules(modulesQ data.ModuleQ, authHeader string) (*resources.EstimatedTimeResponse, error) {
	modules, err := modulesQ.FilterByIsModule(true).Select()
	if err != nil {
		return nil, errors.Wrap(err, "failed to select modules")
	}

	var estimatedTime time.Duration
	for _, module := range modules {
		moduleEstimatedTime, err := helpers.MakeGetEstimatedTimeRequest(data.RequestParams{
			Method: http.MethodPost,
			Link:   module.Link + "/estimate_refresh/module",
			Header: map[string]string{
				"Authorization": authHeader,
			},
			Body:    nil,
			Query:   nil,
			Timeout: 30 * time.Second,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to make refresh request")
		}

		estimatedTime, err = helpers.SumStringDurationWithDuration(moduleEstimatedTime.Data.Attributes.EstimatedTime, estimatedTime)
		if err != nil {
			return nil, errors.Wrap(err, "failed to summarize duration")
		}

		estimatedTime, err = helpers.RoundDuration(estimatedTime.String())
		if err != nil {
			return nil, errors.Wrap(err, "failed to round duration")
		}
	}

	response := NewEstimatedTimeResponse(estimatedTime.String())
	return &response, nil
}

func NewEstimatedTimeModel(estimatedTime string) resources.EstimatedTime {
	return resources.EstimatedTime{
		Key: resources.Key{
			ID:   string(resources.ESTIMATED_TIME),
			Type: resources.ESTIMATED_TIME,
		},
		Attributes: resources.EstimatedTimeAttributes{
			EstimatedTime: estimatedTime,
		},
	}
}

func NewEstimatedTimeResponse(estimatedTime string) resources.EstimatedTimeResponse {
	return resources.EstimatedTimeResponse{
		Data: NewEstimatedTimeModel(estimatedTime),
	}
}
