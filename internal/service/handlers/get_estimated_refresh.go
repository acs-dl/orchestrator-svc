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
	request, err := requests.NewGetEstimatedRefreshRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if request.ModuleName == nil {
		estimatedTimeResponse, err := getEstimatedRefreshModules(helpers.ModulesQ(r), r.Header.Get("Authorization"))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to refresh modules")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, estimatedTimeResponse)
		return
	}

	if request.Submodule == nil {
		estimatedTimeResponse, err := getEstimatedRefreshModule(helpers.ModulesQ(r), *request.ModuleName, r.Header.Get("Authorization"))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to refresh module")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, estimatedTimeResponse)
		return
	}

	estimatedTimeResponse, err := getEstimatedRefreshModuleSubmodules(helpers.ModulesQ(r), *request.ModuleName, r.Header.Get("Authorization"), *request.Submodule)
	if err != nil {
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to refresh module submodule")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	ape.Render(w, estimatedTimeResponse)
}

func getEstimatedRefreshModuleSubmodules(modulesQ data.ModuleQ, moduleName, authHeader, submodules string) (*resources.EstimatedTimeResponse, error) {
	module, err := modulesQ.FilterByNames(moduleName).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get module")
	}

	if module == nil {
		return nil, errors.New("no such module")
	}

	estimatedTime, err := helpers.MakeGetEstimatedTimeRequest(data.RequestParams{
		Method:     http.MethodGet,
		Link:       module.Link + "/refresh/submodule",
		AuthHeader: &authHeader,
		Body:       nil,
		Query: map[string]string{
			"filter[submodule]": submodules,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make refresh request")
	}

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
		Method:     http.MethodGet,
		Link:       module.Link + "/refresh/module",
		AuthHeader: &authHeader,
		Body:       nil,
		Query:      nil,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to make refresh request")
	}

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
			Method:     http.MethodGet,
			Link:       module.Link + "/refresh/module",
			AuthHeader: &authHeader,
			Body:       nil,
			Query:      nil,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to make refresh request")
		}

		estimatedTime, err = helpers.SumStringDurationWithDuration(moduleEstimatedTime.Data.Attributes.EstimatedTime, estimatedTime)
		if err != nil {
			return nil, errors.Wrap(err, "failed to summarize duration")
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
