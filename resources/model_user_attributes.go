/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UserAttributes struct {
	// submodule access level
	AccessLevel string `json:"access_level"`
	// module name
	Module string `json:"module"`
	// phone from module
	Phone *string `json:"phone,omitempty"`
	// submodule name
	Submodule string `json:"submodule"`
	// user id from identity module
	UserId int64 `json:"user_id"`
	// username from module
	Username *string `json:"username,omitempty"`
}
