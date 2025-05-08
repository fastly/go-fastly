package computeacls

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v10/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// ComputeACLID is an ACL Identifier (required).
	ComputeACLID *string
}

// DeleteComputeACL deletes the specified compute ACL.
func Delete(c *fastly.Client, i *DeleteInput) error {
	if i.ComputeACLID == nil {
		return fastly.ErrMissingComputeACLID
	}

	path := fastly.ToSafeURL("resources", "acls", *i.ComputeACLID)

	resp, err := c.Delete(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return err
	}
	defer fastly.CheckCloseForErr(resp.Body.Close)

	if resp.StatusCode != http.StatusOK {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
