package slack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// UpdateConfig is the config object for integration type slack.
type UpdateConfig struct {
	// Webhook is the Slack webhook (required).
	Webhook *string `json:"webhook"`
}

// UpdateInput specifies the information needed for the Update() function to perform
// the operation.
type UpdateInput struct {
	// AlertID is The unique identifier of the workspace alert (required).
	AlertID *string `json:"-"`
	// Config is the configuration associated with the workspace integration (required).
	Config *UpdateConfig `json:"config"`
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Events is a list of event types (required).
	Events *[]string `json:"events"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Update updates the specified slack alert.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Alert, error) {
	if i.AlertID == nil {
		return nil, fastly.ErrMissingAlertID
	}

	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	// Get the current alert to validate the integration type.
	currentAlert, err := Get(ctx, c, &GetInput{
		AlertID:     i.AlertID,
		WorkspaceID: i.WorkspaceID,
		Context:     i.Context,
	})
	if err != nil {
		return nil, err
	}

	// Validate that this is a slack integration
	if currentAlert.Type != IntegrationType {
		return nil, fastly.ErrInvalidType
	}

	if i.Config == nil {
		return nil, fastly.ErrMissingConfig
	}

	// Validate slack integration configuration
	if i.Config.Webhook == nil {
		return nil, fastly.ErrMissingWebhook
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts", *i.AlertID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
