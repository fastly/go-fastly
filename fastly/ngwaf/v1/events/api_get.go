package events

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
	// EventID is the event identifier (required).
	EventID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Get retrieves the specified workspace.
func Get(c *fastly.Client, i *GetInput) (*Event, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.EventID == nil {
		return nil, fastly.ErrMissingEventID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "events", *i.EventID)

	resp, err := c.Get(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var event *Event
	if err := json.NewDecoder(resp.Body).Decode(&event); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return event, nil
}
