package datadog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// CreateInput specifies the information needed for the Create() function to perform
// the operation.
type CreateInput struct {
	// Config is the configuration associated with the workspace integration (required).
	Config Config `json:"-"`
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Events is a list of event types (required).
	Events *string `json:"-"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Create creates a new workspace alert.
func Create(c *fastly.Client, i *CreateInput) (*WorkspaceAlert, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	if (i.Config == Config{}) {
		return nil, fastly.ErrMissingConfig
	}

	// Validate datadog integration configuration
	if i.Config.Key == nil {
		return nil, fastly.ErrMissingKey
	}
	if i.Config.Site == nil {
		return nil, fastly.ErrMissingSite
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts")

	resp, err := c.PostJSON(path, i, fastly.CreateRequestOptions(i.Context))
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