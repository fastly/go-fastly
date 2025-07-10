package virtualpatches

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
	// VirtualPatchID is the virtual patch identifier (required).
	VirtualPatchID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Get retrieves the specified virtual patch.
func Get(c *fastly.Client, i *GetInput) (*VirtualPatch, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	if i.VirtualPatchID == nil {
		return nil, fastly.ErrMissingVirtualPatchID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "virtual-patches", *i.VirtualPatchID)

	resp, err := c.Get(path, fastly.CreateRequestOptions(i.Context))
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
