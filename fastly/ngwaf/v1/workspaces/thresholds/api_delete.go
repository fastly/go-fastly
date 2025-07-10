package thresholds

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
	// ThresholdID is the threshold identifier. Required.
	ThresholdID *string
	// WorkspaceID is the workspace identifier. Required.
	WorkspaceID *string
}

// Delete deletes the specified threshold.
func Delete(c *fastly.Client, i *DeleteInput) error {
	if i.WorkspaceID == nil {
		return fastly.ErrMissingWorkspaceID
	}
	if i.ThresholdID == nil {
		return fastly.ErrMissingThresholdID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "thresholds", *i.ThresholdID)

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
