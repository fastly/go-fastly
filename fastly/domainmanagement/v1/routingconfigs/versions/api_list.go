package versions

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
	// Cursor is the cursor value from the next_cursor field of a previous
	// response, used to retrieve the next page. To request the first page, this
	// should be an empty string or nil.
	Cursor *string
	// Limit is the maximum number of results to return.
	Limit *int
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
	// Sort is the order in which to list the results.
	Sort *string
}

// List retrieves a list of the specified routing config's inactive versions,
// with optional pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Collection, error) {
	if i.RoutingConfigID == nil {
		return nil, fastly.ErrMissingRoutingConfigID
	}

	requestOptions := fastly.CreateRequestOptions()
	if i.Cursor != nil {
		requestOptions.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Sort != nil {
		requestOptions.Params["sort"] = *i.Sort
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "versions")

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
