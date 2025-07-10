package slack

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v10/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to perform
// the operation.
type DeleteInput struct {
	// AlertID is The unique identifier of the workspace alert (required).
	AlertID *string
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Delete deletes the workspace alert.
func Delete(c *fastly.Client, i *DeleteInput) error {
	if i.WorkspaceID == nil {
		return fastly.ErrMissingWorkspaceID
	}
	if i.AlertID == nil {
		return fastly.ErrMissingAlertID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts", *i.AlertID)

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
