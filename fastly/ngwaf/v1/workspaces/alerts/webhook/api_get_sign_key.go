package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// GetKeyInput specifies the information needed for the get signing key
// operation.
type GetKeyInput struct {
	// AlertID is The unique identifier of the workspace alert (required).
	AlertID *string
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// GetKey gets the webhook alert signing key.
func GetKey(ctx context.Context, c *fastly.Client, i *GetKeyInput) (*AlertsKey, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	if i.AlertID == nil {
		return nil, fastly.ErrMissingAlertID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts", *i.AlertID, "signing-key")

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wak *AlertsKey
	if err := json.NewDecoder(resp.Body).Decode(&wak); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return wak, nil
}
