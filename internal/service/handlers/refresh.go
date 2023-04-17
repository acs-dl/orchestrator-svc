package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/requests"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewRefreshRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if request.Data.Attributes.ModuleName == nil {
		err = refreshModules(helpers.ModulesQ(r), r.Header.Get("Authorization"))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to refresh modules")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, http.StatusAccepted)
		return
	}

	if request.Data.Attributes.Submodule == nil {
		err = refreshModule(helpers.ModulesQ(r), *request.Data.Attributes.ModuleName, r.Header.Get("Authorization"))
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to refresh module")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, http.StatusAccepted)
		return
	}

	err = refreshModuleSubmodules(helpers.ModulesQ(r), *request.Data.Attributes.ModuleName, r.Header.Get("Authorization"), *request.Data.Attributes.Submodule)
	if err != nil {
		if err != nil {
			helpers.Log(r).WithError(err).Error("failed to refresh module submodule")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	ape.Render(w, http.StatusAccepted)
}

func refreshModuleSubmodules(modulesQ data.ModuleQ, moduleName, authHeader string, submodules []string) error {
	module, err := modulesQ.FilterByNames(moduleName).Get()
	if err != nil {
		return errors.Wrap(err, "failed to get module")
	}

	if module == nil {
		return errors.New("no such module")
	}

	body := resources.Submodule{
		Key: resources.Key{
			ID:   string(resources.LINKS),
			Type: resources.LINKS,
		},
		Attributes: resources.SubmoduleAttributes{
			Links: submodules,
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "failed to marshal body")
	}

	err = helpers.MakeNoResponseRequest(data.RequestParams{
		Method:     http.MethodPost,
		Link:       module.Link + "/refresh/submodule",
		AuthHeader: &authHeader,
		Body:       bodyBytes,
		Query:      nil,
	})
	if err != nil {
		return errors.Wrap(err, "failed to make refresh request")
	}

	return nil
}

func refreshModule(modulesQ data.ModuleQ, moduleName, authHeader string) error {
	module, err := modulesQ.FilterByNames(moduleName).Get()
	if err != nil {
		return errors.Wrap(err, "failed to get module")
	}

	if module == nil {
		return errors.New("no such module")
	}

	err = helpers.MakeNoResponseRequest(data.RequestParams{
		Method:     http.MethodPost,
		Link:       module.Link + "/refresh/module",
		AuthHeader: &authHeader,
		Body:       nil,
		Query:      nil,
	})
	if err != nil {
		return errors.Wrap(err, "failed to make refresh request")
	}

	return nil
}

func refreshModules(modulesQ data.ModuleQ, authHeader string) error {
	modules, err := modulesQ.FilterByIsModule(true).Select()
	if err != nil {
		return errors.Wrap(err, "failed to select modules")
	}

	for _, module := range modules {
		err = helpers.MakeNoResponseRequest(data.RequestParams{
			Method:     http.MethodPost,
			Link:       module.Link + "/refresh/module",
			AuthHeader: &authHeader,
			Body:       nil,
			Query:      nil,
		})
		if err != nil {
			return errors.Wrap(err, "failed to make refresh request")
		}
	}

	return nil
}
