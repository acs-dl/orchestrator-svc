/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type GetRequest struct {
	// Error during request processing
	Error *string `json:"error,omitempty"`
	// Module to grant permission
	Module string `json:"module"`
	// Already built payload to grant permission
	Payload json.RawMessage `json:"payload"`
	// Status of the request
	Status string `json:"status"`
}