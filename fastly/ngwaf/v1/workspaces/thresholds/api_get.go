package thresholds

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// GetInput specifies the information needed for the Get() function to
// perform the operation.
type GetInput struct {
	// ThresholdID is the threshold identifier (required).
	ThresholdID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Get retrieves the specified threshold.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*Threshold, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.ThresholdID == nil {
		return nil, fastly.ErrMissingThresholdID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "thresholds", *i.ThresholdID)

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
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
