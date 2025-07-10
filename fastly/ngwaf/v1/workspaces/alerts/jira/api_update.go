package jira

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// UpdateConfig is the config object for integration type jira.
type UpdateConfig struct {
	// Host is the name of the Jira instance (required).
	Host *string `json:"host"`
	// IssueType is the Jira issue type associated with the ticket (optional).
	IssueType *string `json:"issue_type,omitempty"`
	// Key is the Jira API key / secret field (required).
	Key *string `json:"key"`
	// Project specifies the Jira project where the issue will be created (required).
	Project *string `json:"project"`
	// Username is the Jira username of the user who created the ticket (required).
	Username *string `json:"username"`
}

// UpdateInput specifies the information needed for the Update() function to perform
// the operation.
type UpdateInput struct {
	// AlertID is The unique identifier of the workspace alert (required).
	AlertID *string `json:"-"`
	// Config is the configuration associated with the workspace integration (required).
	Config UpdateConfig `json:"config"`
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Events is a list of event types (required).
	Events []string `json:"events"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Update updates the specified jira alert.
func Update(c *fastly.Client, i *UpdateInput) (*Alert, error) {
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

	// Validate that this is a jira integration
	if currentAlert.Type != IntegrationType {
		return nil, fastly.ErrInvalidConfigType
	}

	if (i.Config == UpdateConfig{}) {
		return nil, fastly.ErrMissingConfig
	}

	// Validate jira integration configuration
	if i.Config.Host == nil {
		return nil, fastly.ErrMissingHost
	}
	if i.Config.Key == nil {
		return nil, fastly.ErrMissingKey
	}
	if i.Config.Project == nil {
		return nil, fastly.ErrMissingProject
	}
	if i.Config.Username == nil {
		return nil, fastly.ErrMissingUsername
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts", *i.AlertID)

	resp, err := c.PatchJSON(path, i, fastly.CreateRequestOptions(i.Context))
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
