/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type CreateRequest struct {
	// user's id who send request
	FromUser string `json:"from_user"`
	// Module to grant permission
	Module string `json:"module"`
	// Already built payload to grant permission
	Payload json.RawMessage `json:"payload"`
	// user's id for who request was sent
	ToUser string `json:"to_user"`
}
