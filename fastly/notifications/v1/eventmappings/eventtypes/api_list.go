package eventtypes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// ScopeType filters results to event types compatible with the given
	// scope type.
	ScopeType *string
	// Sort orders results alphabetically by event type: "event_type" or
	// "-event_type". Defaults to "event_type".
	Sort *string
}

// List retrieves all audit event types that can be used when creating an
// event mapping.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Collection, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.ScopeType != nil {
		requestOptions.Params["scope_type"] = *i.ScopeType
	}
	if i.Sort != nil {
		requestOptions.Params["sort"] = *i.Sort
	}

	path := fastly.ToSafeURL("notifications", "v1", "event-mappings", "event-types")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cl *Collection
	if err := json.NewDecoder(resp.Body).Decode(&cl); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return cl, nil
}
