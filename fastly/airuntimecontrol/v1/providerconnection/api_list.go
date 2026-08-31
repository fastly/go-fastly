package providerconnection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// Sort is the sort field. Prefix with "-" for descending order (e.g. "-created_at").
	Sort *string
}

// List retrieves all configured provider connections.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*ProviderConnections, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Sort != nil && *i.Sort != "" {
		requestOptions.Params["sort"] = *i.Sort
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "provider-connections")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pcs *ProviderConnections
	if err := json.NewDecoder(resp.Body).Decode(&pcs); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return pcs, nil
}
