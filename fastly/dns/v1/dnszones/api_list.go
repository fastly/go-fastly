package dnszones

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v14/fastly"
)

// ListInput specifies the information needed to list zones.
type ListInput struct {
	// Cursor is the value from the next_cursor field of a previous
	// response, used to retrieve the next page. To request the first
	// page, this should be empty.
	Cursor *string
	// Limit is how many results are returned.
	Limit *int
	// Name filters the list to return only zones that contain the provided name.
	Name *string
	// Sort is the order in which to list the results.
	Sort *string
}

// List retrieves a paginated list of DNS Zones.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Zones, error) {
	path := fastly.ToSafeURL("dns", "v1", "zones")

	requestOptions := fastly.CreateRequestOptions()
	if i.Cursor != nil {
		requestOptions.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Name != nil {
		requestOptions.Params["name"] = *i.Name
	}
	if i.Sort != nil {
		requestOptions.Params["sort"] = *i.Sort
	}

	resp, err := c.GetJSON(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var zones *Zones
	if err := json.NewDecoder(resp.Body).Decode(&zones); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return zones, nil
}
