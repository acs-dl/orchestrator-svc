package helpers

import "encoding/json"

func CreateRefreshModulePayload() (json.RawMessage, error) {
	refreshModulePayload := struct {
		Action string `json:"action"`
	}{
		Action: "refresh_module",
	}

	return json.Marshal(refreshModulePayload)
}
