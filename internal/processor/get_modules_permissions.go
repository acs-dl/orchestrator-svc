package processor

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/acs-dl/orchestrator-svc/internal/data"
	"github.com/acs-dl/orchestrator-svc/internal/service/helpers"
	"github.com/acs-dl/orchestrator-svc/internal/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) handleGetModulesPermissions(msg types.QueueOutput) error {
	p.log.Infof("started handling get modules permissions")

	modules, err := p.modulesQ.Select()
	if err != nil {
		return errors.Wrap(err, "failed to select modules")
	}

	moduleRoles := make(data.Modules)
	moduleRoles[data.ModuleName] = getOrchestratorPermissions()

	for _, module := range modules {
		res, err := helpers.MakeGetRolesRequest(data.RequestParams{
			Method: http.MethodGet,
			Link:   module.Link + "/user_roles",
			Body:   nil,
			Header: map[string]string{
				"Content-Type": "application/json",
			},
			Query:   nil,
			Timeout: 30 * time.Second,
		})
		if err != nil {
			p.log.WithError(err).Errorf("failed to get user roles from `%s`", module.Link+"/user_roles")
			return errors.Wrap(err, "failed to get user roles")
		}
		moduleRoles[module.Name] = res.Data.Attributes
	}

	moduleRolesJson, err := json.Marshal(data.PermissionsPayload{
		RequestId:         *msg.RequestId,
		Action:            SetModulesPermissionsAction,
		ModulePermissions: moduleRoles,
	})
	if err != nil {
		p.log.WithError(err).Errorf("failed to marshal module permission message")
		return err
	}

	err = p.sender.SendMessageToCustomChannel(p.authTopic, p.sender.BuildPermissionsMessage(*msg.RequestId, moduleRolesJson))
	if err != nil {
		p.log.WithError(err).Errorf("failed to send message to custom channel")
		return err
	}

	p.log.Infof("finished handling get modules permissions")
	return nil
}

func getOrchestratorPermissions() data.Roles {
	return data.Roles{
		"super_admin": "write",
		"admin":       "write",
		"user":        "read",
	}
}
