package computeacls

import (
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v9/fastly"
)

// DescribeInput specifies the information needed for the Describe() function to perform
// the operation.
type DescribeInput struct {
	// ComputeACLID is an ACL Identifier (required).
	ComputeACLID *string
}

// Describe describes a specified compute ACL.
func Describe(c *fastly.Client, i *DescribeInput) (*ComputeACL, error) {
	if i.ComputeACLID == nil {
		return nil, fastly.ErrMissingComputeACLID
	}

	path := fastly.ToSafeURL("resources", "acls", *i.ComputeACLID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var acl *ComputeACL
	if err := json.NewDecoder(resp.Body).Decode(&acl); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return acl, nil
}
