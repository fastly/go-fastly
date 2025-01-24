package v1

import (
	"net/http"

	"github.com/fastly/go-fastly/v9/fastly"
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
	Prefix    *string `json:"prefix"`
	Action    *string `json:"action"`
	Operation *string `json:"op"`
}

// Update updates the specified compute ACl.
func Update(c *fastly.Client, i *UpdateInput) error {
	if i.ComputeACLID == nil {
		return fastly.ErrMissingComputeACLID
	}

	path := fastly.ToSafeURL("resources", "acls", *i.ComputeACLID, "entries")

	resp, err := c.PatchJSON(path, i, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
