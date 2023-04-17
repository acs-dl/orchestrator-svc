package data

import (
	"gitlab.com/distributed_lab/acs/orchestrator/resources"
)

const (
	ModuleName         = "orchestrator"
	RemoveModuleAction = "remove_module"
	DeleteUserAction   = "delete_user"
)

type PermissionsPayload struct {
	RequestId         string  `json:"request_id"`
	Action            string  `json:"action"`
	ModulePermissions Modules `json:"module_permissions"`
	ModuleName        string  `json:"module_name"`
}

type DeleteUserPayload struct {
	Action   string  `json:"action"`
	Username *string `json:"username"`
	Phone    *string `json:"phone"`
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

type RequestParams struct {
	Method     string
	Link       string
	Body       []byte
	Query      map[string]string
	AuthHeader *string
}
