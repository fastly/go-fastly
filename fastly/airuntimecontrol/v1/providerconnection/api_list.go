package providerconnection

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// Cursor is the pagination cursor.
	Cursor *string
	// Limit is the maximum number of results per page.
	Limit *int
	// Sort is the sort field. Prefix with "-" for descending order (e.g. "-created_at").
	Sort *string
}

// List retrieves all configured provider connections for a customer, with
// optional pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*ProviderConnections, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Cursor != nil && *i.Cursor != "" {
		requestOptions.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
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
