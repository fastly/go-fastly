package computeacls

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v12/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// ComputeACLID is an ACL Identifier (required).
	ComputeACLID *string
}

// DeleteComputeACL deletes the specified compute ACL.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.ComputeACLID == nil {
		return fastly.ErrMissingComputeACLID
	}

	path := fastly.ToSafeURL("resources", "acls", *i.ComputeACLID)

	resp, err := c.Delete(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
