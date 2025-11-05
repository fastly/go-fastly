package redactions

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v12/fastly"
)

// DeleteInput specifies the information needed for the Delete()
// function to perform the operation.
type DeleteInput struct {
	// RedactionID is the redaction identifier (required).
	RedactionID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Delete deletes the specified redaction.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.WorkspaceID == nil {
		return fastly.ErrMissingWorkspaceID
	}
	if i.RedactionID == nil {
		return fastly.ErrMissingRedactionID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "redactions", *i.RedactionID)

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
