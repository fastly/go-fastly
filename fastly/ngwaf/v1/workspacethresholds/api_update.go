package workspacethresholds

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// UpdateInput specifies the information needed for the Put() function to perform
// the operation.
type UpdateInput struct {
	// Action
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Action is the action to take when the theshold is exceeded.
	Action string `json:"action"`
	// DontNotify determines whether to silence notifications when action is taken.
	DontNotify bool `json:"dont_notify"`
	// Duration is the duration of the action in place.
	Duration int `json:"duration"`
	// Enabled is whether this threshold is active.
	Enabled bool `json:"enabled"`
	// Interval is the threshold interval in seconds.
	Interval int `json:"interval"`
	// Limit is the threshold limit.
	Limit int `json:"limit"`
	// Name is the threshold name.
	Name string `json:"name"`
	// Signal is the name of the signal this threshold is acting on.
	Signal string `json:"signal"`
	// ThresholdID is the workspace threshold identifier (required).
	ThresholdID *string `json:"-"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Update updates the specified workspace threshold.
func Update(c *fastly.Client, i *UpdateInput) (*UpdateWorkspaceThreshold, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.ThresholdID == nil {
		return nil, fastly.ErrMissingThresholdID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "thresholds", *i.ThresholdID)

	resp, err := c.PatchJSON(path, i, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var uwt *UpdateWorkspaceThreshold
	if err := json.NewDecoder(resp.Body).Decode(&uwt); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return uwt, nil
}
