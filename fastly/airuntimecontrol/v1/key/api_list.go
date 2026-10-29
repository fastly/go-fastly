package key

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v18/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// Model filters results by AI model identifier.
	Model *string
	// Provider filters results by AI model provider.
	Provider *string
	// IncludeDeleted includes deleted virtual keys in the results.
	IncludeDeleted *bool
	// Search filters results to virtual keys whose name contains this substring.
	Search *string
	// Limit limits how many results are returned per page (1-100).
	Limit *int
	// Sort is the sort field. Prefix with `-` for descending order (e.g.,
	// `sort=-created_at`). One of `created_at`, `updated_at`, `name`,
	// `expires_at`, `user_name`, `last_used_at` (each optionally prefixed
	// with `-`).
	Sort *string
}

// List lists the virtual keys belonging to your account. It automatically
// paginates through all pages returned by the API.
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
