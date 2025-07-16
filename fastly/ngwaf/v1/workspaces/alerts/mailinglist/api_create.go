package mailinglist

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// CreateConfig is the config object for integration type mailinglist.
type CreateConfig struct {
	// Address An email address (required).
	Address *string `json:"address"`
}

// CreateInput specifies the information needed for the Create() function to perform
// the operation.
type CreateInput struct {
	// Config is the configuration associated with the workspace integration (required).
	Config *CreateConfig `json:"config"`
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Description is an optional description for the alert.
	Description *string `json:"description,omitempty"`
	// Events is a list of event types (required).
	Events *[]string `json:"events"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Private struct to ensure correct alert type.
type privateMailingListInput struct {
	Config      *CreateConfig    `json:"config"`
	Context     *context.Context `json:"-"`
	Description *string          `json:"description,omitempty"`
	Events      *[]string        `json:"events"`
	Type        *string          `json:"type"`
	WorkspaceID *string          `json:"-"`
}

// Create creates a new mailinglist alert.
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
	// Validate mailinglist integration configuration
	if i.Config.Address == nil {
		return nil, fastly.ErrMissingAddress
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts")

	mailingListInput := privateMailingListInput{
		Config:      i.Config,
		Context:     i.Context,
		Description: i.Description,
		Events:      i.Events,
		Type:        fastly.ToPointer(IntegrationType),
		WorkspaceID: i.WorkspaceID,
	}

	resp, err := c.PostJSON(ctx, path, mailingListInput, fastly.CreateRequestOptions())
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
