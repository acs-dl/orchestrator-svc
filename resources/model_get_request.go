/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type GetRequest struct {
	// timestamp when request was created
	CreatedAt string `json:"created_at"`
	// Error during request processing
	Error *string `json:"error,omitempty"`
	// user's id who send request
	FromUser string `json:"from_user"`
	// Module to grant permission
	Module string `json:"module"`
	// Already built payload to grant permission
	Payload json.RawMessage `json:"payload"`
	// Status of the request
	Status string `json:"status"`
	// user's id for who request was sent
	ToUser string `json:"to_user"`
}
