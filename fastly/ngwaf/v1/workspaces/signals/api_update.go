package signals

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// Description is the new description for the signal (required).
	Description *string `json:"description"`
	// SignalID is the id of the signal that's being updated (required).
	SignalID *string `json:"-"`
	// WorkspaceID is the ID of the workspace that the signal belongs to.
	WorkspaceID *string `json:"-"`
}

// Update updates the specified signal.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Signal, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.SignalID == nil {
		return nil, fastly.ErrMissingSignalID
	}
	if i.Description == nil {
		return nil, fastly.ErrMissingDescription
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "signals", *i.SignalID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var signal *Signal
	if err := json.NewDecoder(resp.Body).Decode(&signal); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}
	return signal, nil
}
