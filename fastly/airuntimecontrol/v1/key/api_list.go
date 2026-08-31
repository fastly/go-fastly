package key

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
	// Model filters results by AI model identifier.
	Model *string
	// Provider filters results by AI model provider.
	Provider *string
	// IncludeDeleted includes deleted virtual keys when true.
	IncludeDeleted *bool
	// Search filters virtual keys by substring match on key name.
	Search *string
	// Limit is the maximum number of results per page.
	Limit *int
	// Sort is the sort field. Prefix with "-" for descending order (e.g. "-created_at").
	Sort *string
}

// List retrieves all virtual keys, automatically paginating through all
// pages returned by the API, with optional filtering.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]VirtualKeyListItem, error) {
	var (
		out    []VirtualKeyListItem
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

// listPage retrieves a single page of virtual keys.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, cursor *string) (*VirtualKeys, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Model != nil && *i.Model != "" {
		requestOptions.Params["model"] = *i.Model
	}
	if i.Provider != nil && *i.Provider != "" {
		requestOptions.Params["provider"] = *i.Provider
	}
	if i.IncludeDeleted != nil {
		requestOptions.Params["include_deleted"] = strconv.FormatBool(*i.IncludeDeleted)
	}
	if i.Search != nil && *i.Search != "" {
		requestOptions.Params["search"] = *i.Search
	}
	if cursor != nil && *cursor != "" {
		requestOptions.Params["cursor"] = *cursor
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Sort != nil && *i.Sort != "" {
		requestOptions.Params["sort"] = *i.Sort
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "keys")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var keys *VirtualKeys
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return keys, nil
}
