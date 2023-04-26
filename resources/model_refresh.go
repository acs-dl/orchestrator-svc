/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Refresh struct {
	Key
	Attributes RefreshAttributes `json:"attributes"`
}
type RefreshResponse struct {
	Data     Refresh  `json:"data"`
	Included Included `json:"included"`
}

type RefreshListResponse struct {
	Data     []Refresh `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustRefresh - returns Refresh from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRefresh(key Key) *Refresh {
	var refresh Refresh
	if c.tryFindEntry(key, &refresh) {
		return &refresh
	}
	return nil
}
