/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ModuleInfo struct {
	Key
	Attributes ModuleInfoAttributes `json:"attributes"`
}
type ModuleInfoResponse struct {
	Data     ModuleInfo `json:"data"`
	Included Included   `json:"included"`
}

type ModuleInfoListResponse struct {
	Data     []ModuleInfo `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustModuleInfo - returns ModuleInfo from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustModuleInfo(key Key) *ModuleInfo {
	var moduleInfo ModuleInfo
	if c.tryFindEntry(key, &moduleInfo) {
		return &moduleInfo
	}
	return nil
}
