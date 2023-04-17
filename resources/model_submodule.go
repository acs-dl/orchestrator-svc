/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Submodule struct {
	Key
	Attributes SubmoduleAttributes `json:"attributes"`
}
type SubmoduleResponse struct {
	Data     Submodule `json:"data"`
	Included Included  `json:"included"`
}

type SubmoduleListResponse struct {
	Data     []Submodule `json:"data"`
	Included Included    `json:"included"`
	Links    *Links      `json:"links"`
}

// MustSubmodule - returns Submodule from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSubmodule(key Key) *Submodule {
	var submodule Submodule
	if c.tryFindEntry(key, &submodule) {
		return &submodule
	}
	return nil
}
