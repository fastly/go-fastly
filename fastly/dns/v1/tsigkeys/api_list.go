package tsigkeys

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v14/fastly"
)

// ListInput specifies the information needed to list TSIG keys.
type ListInput struct {
	// Limit is how many results are returned per page.
	Limit *int
	// Name filters the list to return only TSIG keys that contain the provided name.
	Name *string
	// Sort is the order in which to list the results.
	Sort *string
}

// List retrieves all TSIG keys, automatically paginating through all pages.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]TSIGKey, error) {
	var (
		out    []TSIGKey
		cursor *string
	)
	for {
		page, err := listPage(ctx, c, i, cursor)
		if err != nil {
			return nil, err
		}
		out = append(out, page.Data...)
		if page.Meta.NextCursor == nil || *page.Meta.NextCursor == "" {
			break
		}
		cursor = page.Meta.NextCursor
	}
	return out, nil
}

// listPage retrieves a single page of TSIG keys.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, cursor *string) (*TSIGKeys, error) {
	path := fastly.ToSafeURL("dns", "v1", "tsig-keys")

	requestOptions := fastly.CreateRequestOptions()
	if cursor != nil {
		requestOptions.Params["cursor"] = *cursor
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

	var tsigKeys *TSIGKeys
	if err := json.NewDecoder(resp.Body).Decode(&tsigKeys); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return tsigKeys, nil
}
