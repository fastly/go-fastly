package computeacls

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fastly/go-fastly/v9/fastly"
)

// LookupInput specifies the information needed for the Lookup() function to perform
// the operation.
type LookupInput struct {
	// ComputeACLID is an ACL Identifier (required).
	ComputeACLID *string
	// ComputeACLIP is a valid IPv4 or IPv6 address (required).
	ComputeACLIP *string
}

// Lookup finds a matching ACL entry for an IP address.
func Lookup(c *fastly.Client, i *LookupInput) (*ComputeACLEntry, error) {
	if i.ComputeACLID == nil {
		return nil, fastly.ErrMissingComputeACLID
	}
	if i.ComputeACLIP == nil {
		return nil, fastly.ErrMissingComputeACLIP
	}

	path := fastly.ToSafeURL("resources", "acls", *i.ComputeACLID, "entry", *i.ComputeACLIP)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var entry *ComputeACLEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return entry, nil
}
