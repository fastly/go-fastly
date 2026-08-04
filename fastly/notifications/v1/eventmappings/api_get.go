package eventmappings

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetInput specifies the information needed for the Get() function to
// perform the operation.
type GetInput struct {
	// MappingID is the ID of the event mapping (required).
	MappingID *string
}

// Get retrieves the specified event mapping.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*EventMapping, error) {
	if i.MappingID == nil {
		return nil, fastly.ErrMissingMappingID
	}

	path := fastly.ToSafeURL("notifications", "v1", "event-mappings", *i.MappingID)

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var em *EventMapping
	if err := json.NewDecoder(resp.Body).Decode(&em); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return em, nil
}
