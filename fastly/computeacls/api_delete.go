package computeacls

import (
	"net/http"

	"github.com/fastly/go-fastly/v9/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// ComputeACLID is an ACL Identifier (required).
	ComputeACLID *string
}

// DeleteComputeACL deletes the specified compute ACL.
func Delete(c *fastly.Client, i *DeleteInput) error {
	if i.ComputeACLID == nil {
		return fastly.ErrMissingComputeACLID
	}

	path := fastly.ToSafeURL("resources", "acls", *i.ComputeACLID)

	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
