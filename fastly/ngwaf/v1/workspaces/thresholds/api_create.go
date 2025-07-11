package thresholds

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Action to take when threshold is exceeded. Must be one of `block` or `log`. Required.
	Action *string `json:"action"`
	// DontNotify indicates whether to silence notifications when action is taken. Defaults to false.
	DontNotify *bool `json:"dont_notify,omitempty"`
	// Duration is the duration the action is in place. Default duration is 0. If set, must be greater than 0 and less than 31556900 (1 year).
	Duration *int `json:"duration,omitempty"`
	// Enabled is whether this threshold is active. Defaults to false.
	Enabled *bool `json:"enabled,omitempty"`
	// Limit is the threshold limit. Must be between 1 and 10000 inclusive. Required.
	Limit *int `json:"limit"`
	// Interval is the threshold interval in seconds. Must be one of 60, 600, or 36000. Required.
	Interval *int `json:"interval"`
	// Name is the threshold name. Required.
	Name *string `json:"name"`
	// Signal is the name of the signal this threshold is acting on. Required.
	Signal *string `json:"signal"`
	// WorkspaceID is the ID of the workspace that the threshold is being created in. Required.
	WorkspaceID *string `json:"-"`
}

// Create creates a new threshold.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Threshold, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if i.Action == nil {
		return nil, fastly.ErrMissingAction
	}
	if i.Limit == nil {
		return nil, fastly.ErrMissingLimit
	}
	if i.Interval == nil {
		return nil, fastly.ErrMissingInterval
	}
	if i.Signal == nil {
		return nil, fastly.ErrMissingSignal
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "thresholds")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var threshold *Threshold
	if err := json.NewDecoder(resp.Body).Decode(&threshold); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return threshold, nil
}
