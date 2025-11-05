package thresholds

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v12/fastly"
)

// UpdateInput specifies the information needed for the Update()
// function to perform the operation.
type UpdateInput struct {
	// Action to take when threshold is exceeded. Must be one of
	// `block` or `log` (required).
	Action *string `json:"action"`
	// DontNotify indicates whether to silence notifications when
	// action is taken.
	DontNotify *bool `json:"dont_notify,omitempty"`
	// Duration is the duration the action is in place. Must be
	// greater than 0 and less than 31556900 (1 year).
	Duration *int `json:"duration,omitempty"`
	// Enabled is whether this threshold is active.
	Enabled *bool `json:"enabled,omitempty"`
	// Limit is the threshold limit. Must be between 1 and 10000
	// inclusive.
	Limit *int `json:"limit,omitempty"`
	// Interval is the threshold interval in seconds. Must be one
	// of 60, 600, or 36000.
	Interval *int `json:"interval,omitempty"`
	// Name is the threshold name.
	Name *string `json:"name,omitempty"`
	// Signal is the name of the signal this threshold is acting
	// on.
	Signal *string `json:"signal,omitempty"`
	// ThresholdID is the threshold identifier (required).
	ThresholdID *string `json:"-"`
	// WorkspaceID is the ID of the workspace that the threshold
	// being updated is in (required).
	WorkspaceID *string `json:"-"`
}

// Update updates the specified workspace.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Threshold, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.ThresholdID == nil {
		return nil, fastly.ErrMissingThresholdID
	}
	if i.Action == nil {
		return nil, fastly.ErrMissingAction
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "thresholds", *i.ThresholdID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
