package routingconfigs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

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
	// Sort is the order in which to list the results.
	Sort *string
	// State filters results by lifecycle state. If multiple values are
	// provided, they will be sent as a comma-separated string.
	State []string
}

// List retrieves a list of routing configs, with optional filtering and
// pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Collection, error) {
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
	if len(i.State) > 0 {
		requestOptions.Params["state"] = strings.Join(i.State, ",")
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs")

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
