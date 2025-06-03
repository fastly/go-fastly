package virtualpatches

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
	// Enabled is the toggle status indicator of the VirtualPatch (optional).
	Enabled *bool `json:"enabled"`
	// Mode is action to take when a signal for virtual patch is detected (optional).
	Mode *string `json:"mode"`
	// VirtualPatchID is the virtual patch identifier (required).
	VirtualPatchID *string `json:"-"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Update updates the specified virtual patch.
func Update(c *fastly.Client, i *UpdateInput) (*VirtualPatch, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.VirtualPatchID == nil {
		return nil, fastly.ErrMissingVirtualPatchID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "virtual-patches", *i.VirtualPatchID)

	resp, err := c.PatchJSON(path, i, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vp *VirtualPatch
	if err := json.NewDecoder(resp.Body).Decode(&vp); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return vp, nil
}
