/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ModulesRoles struct {
	Key
	Attributes ModulesRolesAttributes `json:"attributes"`
}
type ModulesRolesResponse struct {
	Data     ModulesRoles `json:"data"`
	Included Included     `json:"included"`
}

type ModulesRolesListResponse struct {
	Data     []ModulesRoles `json:"data"`
	Included Included       `json:"included"`
	Links    *Links         `json:"links"`
}

// MustModulesRoles - returns ModulesRoles from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustModulesRoles(key Key) *ModulesRoles {
	var modulesRoles ModulesRoles
	if c.tryFindEntry(key, &modulesRoles) {
		return &modulesRoles
	}
	return nil
}
