package mailinglist

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// Config is the config object for integration type mailinglist.
type Config struct {
	// Address An email address (required).
	Address *string `json:"address,omitempty"`
}

// UpdateInput specifies the information needed for the Update() function to perform
// the operation.
type UpdateInput struct {
	// AlertID is The unique identifier of the workspace alert (required).
	AlertID *string `json:"-"`
	// Config is the configuration associated with the workspace integration (required).
	Config Config `json:"config"`
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Events is a list of event types (required).
	Events *string `json:"events"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Update updates the specified workspace alert.
func Update(c *fastly.Client, i *UpdateInput) (*WorkspaceAlert, error) {
	if i.AlertID == nil {
		return nil, fastly.ErrMissingAlertID
	}

	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	// Get the current alert to validate the integration type.
	currentAlert, err := Get(c, &GetInput{
		AlertID:     i.AlertID,
		WorkspaceID: i.WorkspaceID,
		Context:     i.Context,
	})
	if err != nil {
		return nil, err
	}

	// Validate that this is a mailinglist integration
	if currentAlert.Type != IntegrationType {
		return nil, fastly.ErrInvalidConfigType
	}

	if (i.Config == Config{}) {
		return nil, fastly.ErrMissingConfig
	}

	// Validate mailinglist integration configuration
	if i.Config.Address == nil {
		return nil, fastly.ErrMissingAddress
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts", *i.AlertID)

	resp, err := c.PatchJSON(path, i, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wa *WorkspaceAlert
	if err := json.NewDecoder(resp.Body).Decode(&wa); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return wa, nil
}