package paths

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
	// Limit is the maximum number of results to return per page.
	Limit *int
	// Match filters results using the given path pattern matching strategy.
	Match *string
	// Path filters results by path pattern.
	Path *string
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
	// Sort is the order in which to list the results.
	Sort *string
}

// List retrieves all paths belonging to the specified routing config, with
// optional filtering, automatically paginating through all pages.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]Data, error) {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return nil, fastly.ErrMissingRoutingConfigID
	}

	var (
		out    []Data
		cursor *string
	)
	for {
		page, err := listPage(ctx, c, i, cursor)
		if err != nil {
			return nil, err
		}
		out = append(out, page.Data...)
		if page.Meta.NextCursor == "" {
			break
		}
		cursor = &page.Meta.NextCursor
	}
	return out, nil
}

// listPage retrieves a single page of paths belonging to the specified
// routing config.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, cursor *string) (*Collection, error) {
	requestOptions := fastly.CreateRequestOptions()
	if cursor != nil {
		requestOptions.Params["cursor"] = *cursor
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Sort != nil {
		requestOptions.Params["sort"] = *i.Sort
	}
	if i.Path != nil {
		requestOptions.Params["path"] = *i.Path
	}
	if i.Match != nil {
		requestOptions.Params["match"] = *i.Match
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "paths")

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
