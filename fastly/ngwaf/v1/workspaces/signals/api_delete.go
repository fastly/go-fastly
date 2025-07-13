package signals

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v11/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// SignalID is the signal identifier (required).
	SignalID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Delete deletes the specified signal.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.WorkspaceID == nil {
		return fastly.ErrMissingWorkspaceID
	}
	if i.SignalID == nil {
		return fastly.ErrMissingSignalID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "signals", *i.SignalID)

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
