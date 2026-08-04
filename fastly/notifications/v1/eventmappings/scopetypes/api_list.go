package scopetypes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// Sort orders results alphabetically by scope type: "scope_type" or
	// "-scope_type". Defaults to "scope_type".
	Sort *string
}

// List retrieves all scope types supported by the Event Mappings API.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Collection, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Sort != nil {
		requestOptions.Params["sort"] = *i.Sort
	}

	path := fastly.ToSafeURL("notifications", "v1", "event-mappings", "scope-types")

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
