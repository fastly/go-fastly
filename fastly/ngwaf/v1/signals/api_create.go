package signals

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Description is a description of the signal (optional).
	Description *string `json:"description,omitempty"`
	// Name is the name of the signal. Must be between 3 and 25 characters. Letters, numbers, hyphens, and spaces are accepted. Special characters and periods are not accepted.
	Name *string `json:"name"`
	// WorkspaceID is the ID of the workspace that the signal is being created in.
	WorkspaceID *string `json:"-"`
}

// Create creates a new signal in the given workspace.
func Create(c *fastly.Client, i *CreateInput) (*Signal, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "signals")

	resp, err := c.PostJSON(path, i, fastly.CreateRequestOptions(i.Context))
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
