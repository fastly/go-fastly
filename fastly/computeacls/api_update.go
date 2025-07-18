package computeacls

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v10/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// ComputeACLID  is an ACL Identifier (required).
	ComputeACLID *string
	// Entries is a list of ACL entries.
	Entries []*BatchComputeACLEntry `json:"entries"`
}

// BatchComputeACLEntry is a subset of the UpdateInput input structure.
type BatchComputeACLEntry struct {
	// Prefix is an IP prefix defined in Classless Inter-Domain Routing (CIDR) format.
	Prefix *string `json:"prefix"`
	// Action is one of "ALLOW" or "BLOCK".
	Action *string `json:"action"`
	// Operation is one of "create", "update", or "delete".
	Operation *string `json:"op"`
}

// Update updates the specified compute ACl.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) error {
	if i.ComputeACLID == nil {
		return fastly.ErrMissingComputeACLID
	}

	path := fastly.ToSafeURL("resources", "acls", *i.ComputeACLID, "entries")

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
