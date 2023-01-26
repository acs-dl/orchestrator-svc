/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ModuleIcon struct {
	Key
	Attributes ModuleIconAttributes `json:"attributes"`
}
type ModuleIconResponse struct {
	Data     ModuleIcon `json:"data"`
	Included Included   `json:"included"`
}

type ModuleIconListResponse struct {
	Data     []ModuleIcon `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustModuleIcon - returns ModuleIcon from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustModuleIcon(key Key) *ModuleIcon {
	var moduleIcon ModuleIcon
	if c.tryFindEntry(key, &moduleIcon) {
		return &moduleIcon
	}
	return nil
}
