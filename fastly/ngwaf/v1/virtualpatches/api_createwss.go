package virtualpatches

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// CreateWorkspace specifies the information needed for the Create() function to
// perform the operation.
type CreateWorkspace struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Name is the name of a workspace to create (required).
	Name *string `json:"name"`
	// Description is the description of a workspace.
	Description *string `json:"description"`
	// Mode is the mode of a workspace.
	Mode *string `json:"mode"`
}

// Create creates a new workspace.
func Create(c *fastly.Client, i *CreateWorkspace) (*Workspace, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}

	resp, err := c.PostJSON("/ngwaf/v1/workspaces", i, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ws *Workspace
	if err := json.NewDecoder(resp.Body).Decode(&ws); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return ws, nil
}
