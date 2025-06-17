package workspacethresholds

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// GetInput specifies the information needed for the Get() function to perform
// the operation.
type GetInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// ThresholdID is the workspace threshold identifier (required).
	ThresholdID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Get retrieves the specified workspace threshold.
func Get(c *fastly.Client, i *GetInput) (*WorkspaceThreshold, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	if i.ThresholdID == nil {
		return nil, fastly.ErrMissingThresholdID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "thresholds", *i.ThresholdID)

	resp, err := c.Get(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wt *WorkspaceThreshold
	if err := json.NewDecoder(resp.Body).Decode(&wt); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return wt, nil
}
