package handlers

import (
	"net/http"

	"github.com/acs-dl/orchestrator-svc/internal/data"
	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"github.com/acs-dl/orchestrator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetModules(w http.ResponseWriter, r *http.Request) {
	modules, err := helpers.ModulesQ(r).FilterByIsModule(true).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get modules")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newModuleInfoListResponse(modules))
}

func newModuleInfoListResponse(modules []data.Module) resources.ModuleInfoListResponse {
	return resources.ModuleInfoListResponse{
		Data: newModuleInfoList(modules),
	}
}

func newModuleInfoList(modules []data.Module) []resources.ModuleInfo {
	result := make([]resources.ModuleInfo, 0)

	for i := range modules {
		result = append(result, newModuleInfo(modules[i]))
	}

	return result
}

func newModuleInfo(module data.Module) resources.ModuleInfo {
	return resources.ModuleInfo{
		Key: resources.Key{
			ID:   module.Name,
			Type: resources.MODULES,
		},
		Attributes: resources.ModuleInfoAttributes{
			Icon:   data.ModuleIcon[module.Name],
			Name:   module.Title,
			Prefix: module.Prefix,
		},
	}
}
