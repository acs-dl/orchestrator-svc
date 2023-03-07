package handlers

import (
	"encoding/json"
	"fmt"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/helpers"
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	modules, err := helpers.ModulesQ(r).Select()
	if err != nil {
		helpers.Log(r).WithError(err).Errorf("failed to select modules")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := newModulesRolesResponse()

	if len(modules) == 0 {
		helpers.Log(r).Errorf("no modules was found")
		ape.Render(w, response)
		return
	}

	for _, module := range modules {
		moduleRoles, err := makeGetRolesRequest(module.Link)
		if err != nil {
			helpers.Log(r).WithError(err).Infof("failed to get roles from module `%s`", module.Name)
			ape.RenderErr(w, problems.InternalError())
			return
		}
		response.Data.Attributes[module.Name] = moduleRoles.Data.Attributes
	}

	ape.Render(w, response)
}

func makeGetRolesRequest(moduleLink string) (*ModuleRolesResponse, error) {
	link := fmt.Sprintf(moduleLink + "/roles")
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create request")
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error making http request")
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("error in response, status %s", res.Status))
	}

	var returned ModuleRolesResponse

	if err := json.NewDecoder(res.Body).Decode(&returned); err != nil {
		return nil, errors.Wrap(err, " failed to unmarshal body")
	}

	return &returned, nil
}

func newModulesRolesResponse() ModulesRolesResponse {
	return ModulesRolesResponse{
		Data: ModulesRoles{
			Key: resources.Key{
				ID:   "0",
				Type: resources.MODULES,
			},
			Attributes: Modules{},
		},
	}
}

type ModulesRolesResponse struct {
	Data ModulesRoles `json:"data"`
}

type ModulesRoles struct {
	resources.Key
	Attributes Modules `json:"attributes"`
}

type Modules map[string]Roles

type Roles map[string]string
type ModuleRoles struct {
	resources.Key
	Attributes Roles `json:"attributes"`
}

type ModuleRolesResponse struct {
	Data ModuleRoles `json:"data"`
}
