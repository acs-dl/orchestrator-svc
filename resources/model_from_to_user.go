/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type FromToUser struct {
	Key
	Attributes FromToUserAttributes `json:"attributes"`
}
type FromToUserResponse struct {
	Data     FromToUser `json:"data"`
	Included Included   `json:"included"`
}

type FromToUserListResponse struct {
	Data     []FromToUser `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustFromToUser - returns FromToUser from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustFromToUser(key Key) *FromToUser {
	var fromToUser FromToUser
	if c.tryFindEntry(key, &fromToUser) {
		return &fromToUser
	}
	return nil
}
