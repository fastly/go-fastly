package datadog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v11/fastly"
)

// Config is the config object for integration type datadog.
type CreateConfig struct {
	// Key is the Datadog integration key (required).
	Key *string `json:"key"`
	// Site is the Datadog site (required).
	Site *string `json:"site"`
}

// CreateInput specifies the information needed for the Create() function to perform
// the operation.
type CreateInput struct {
	// Config is the configuration associated with the workspace integration (required).
	Config *CreateConfig
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Description is an optional description for the alert.
	Description *string
	// Events is a list of event types (required).
	Events *[]string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Private struct to ensure correct alert type.
type privateDatadogInput struct {
	Config      *CreateConfig    `json:"config"`
	Context     *context.Context `json:"-"`
	Description *string          `json:"description,omitempty"`
	Events      *[]string        `json:"events"`
	Type        *string          `json:"type"`
	WorkspaceID *string          `json:"-"`
}

// Create creates a new datadog alert.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Alert, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.Config == nil {
		return nil, fastly.ErrMissingConfig
	}
	if i.Events == nil {
		return nil, fastly.ErrMissingEvents
	}

	// Validate datadog integration configuration
	if i.Config.Key == nil {
		return nil, fastly.ErrMissingKey
	}
	if i.Config.Site == nil {
		return nil, fastly.ErrMissingSite
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts")

	datadogInput := privateDatadogInput{
		Config:      i.Config,
		Context:     i.Context,
		Description: i.Description,
		Events:      i.Events,
		Type:        fastly.ToPointer(IntegrationType),
		WorkspaceID: i.WorkspaceID,
	}

	resp, err := c.PostJSON(ctx, path, datadogInput, fastly.CreateRequestOptions())
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
