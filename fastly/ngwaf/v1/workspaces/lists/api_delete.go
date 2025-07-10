package lists

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
	// ListID is the id of the list to be deleted (required).
	ListID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Delete deletes the specified list.
func Delete(c *fastly.Client, i *DeleteInput) error {
	if i.WorkspaceID == nil {
		return fastly.ErrMissingWorkspaceID
	}
	if i.ListID == nil {
		return fastly.ErrMissingListID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "lists", *i.ListID)

	resp, err := c.Delete(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
