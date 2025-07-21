package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v11/fastly"
)

// GetInput specifies the information needed for the Get() function to perform
// the operation.
type GetInput struct {
	// AlertID is The unique identifier of the workspace alert (required).
	AlertID *string
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Get retrieves the webhook workspace alert.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*Alert, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.AlertID == nil {
		return nil, fastly.ErrMissingAlertID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts", *i.AlertID)

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wa *Alert
	if err := json.NewDecoder(resp.Body).Decode(&wa); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return wa, nil
}
