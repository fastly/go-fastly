package workspaces

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v11/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Delete deletes the specified workspace.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.WorkspaceID == nil {
		return fastly.ErrMissingWorkspaceID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID)

	resp, err := c.Delete(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
