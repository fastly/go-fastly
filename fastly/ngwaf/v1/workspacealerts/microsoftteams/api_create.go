package microsoftteams

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// CreateConfig is the config object for integration type microsoftteams.
type CreateConfig struct {
	// Webhook is the Microsoft Teams webhook (required).
	Webhook *string `json:"webhook"`
}

// CreateInput specifies the information needed for the Create() function to perform
// the operation.
type CreateInput struct {
	// Config is the configuration associated with the workspace integration (required).
	Config CreateConfig `json:"config"`
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Description is an optional description for the alert (optional).
	Description *string `json:"description"`
	// Events is a list of event types (required).
	Events *string `json:"events"`
	// Type is the type of the workspace integration (required).
	Type *string `json:"type"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Create creates a new workspace alert.
func Create(c *fastly.Client, i *CreateInput) (*WorkspaceAlert, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.Type == nil || *i.Type != IntegrationType {
		return nil, fastly.ErrInvalidConfigType
	}
	if (i.Config == CreateConfig{}) {
		return nil, fastly.ErrMissingConfig
	}
	if i.Events == nil {
		return nil, fastly.ErrMissingEvents
	}

	// Validate microsoftteams integration configuration
	if i.Config.Webhook == nil {
		return nil, fastly.ErrMissingWebhook
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