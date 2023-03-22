package processor

import (
	"encoding/json"

	"gitlab.com/distributed_lab/acs/orchestrator/internal/data"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/service/handlers"
	"gitlab.com/distributed_lab/acs/orchestrator/internal/types"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (p *processor) handleGetModulesPermissions(msg types.QueueOutput) error {
	modules, err := p.modulesQ.Select()
	if err != nil {
		return errors.Wrap(err, "failed to select modules")
	}

	moduleRoles := make(data.Modules)
	moduleRoles[data.ModuleName] = getOrchestratorPermissions()

	for _, module := range modules {
		res, err := handlers.MakeGetRolesRequest(module.Link, "/user_roles")
		if err != nil {
			p.log.WithError(err).Errorf("failed to get user roles from `%s`", module)
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

	err = p.sender.SendMessageToCustomChannel("auth", p.sender.BuildPermissionsMessage(*msg.RequestId, moduleRolesJson))
	if err != nil {
		p.log.WithError(err).Errorf("failed to send message to custom channel")
		return err
	}

	return nil
}

func getOrchestratorPermissions() data.Roles {
	return data.Roles{
		"super_admin": "write",
		"admin":       "write",
		"user":        "read",
	}
}
