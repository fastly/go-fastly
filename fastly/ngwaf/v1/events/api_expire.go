package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// ExpireInput specifies the information needed for the Expire() function to
// perform the operation.
type ExpireInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// EventID is the event identifier (required).
	EventID *string `json:"-"`
	// IsExpired sets the value of IsExpired on the event (required).
	IsExpired *bool `json:"is_expired"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Expire expires the specified event.
func Expire(c *fastly.Client, i *ExpireInput) (*Event, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.EventID == nil {
		return nil, fastly.ErrMissingEventID
	}
	if i.IsExpired == nil {
		return nil, fastly.ErrMissingIsExpired
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "events", *i.EventID)

	resp, err := c.PatchJSON(path, i, fastly.CreateRequestOptions(i.Context))
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
