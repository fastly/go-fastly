package session

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// Key filters results by a specific virtual key ID.
	Key *string
	// Provider filters results by provider name.
	Provider *string
	// Model filters results by model name.
	Model *string
	// From is the (inclusive) start of the time range.
	From *time.Time
	// To is the (inclusive) end of the time range. Defaults to now.
	To *time.Time
	// Limit is the maximum number of results per page.
	Limit *int
	// Sort is the sort field. Prefix with "-" for descending order (e.g. "-created_at").
	Sort *string
}

// List retrieves session logs including request/response data and metadata,
// automatically paginating through all pages, with optional filtering.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]Session, error) {
	var (
		out    []Session
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

// listPage retrieves a single page of session logs.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, cursor *string) (*Sessions, error) {
	requestOptions := fastly.CreateRequestOptions()
	if i.Key != nil && *i.Key != "" {
		requestOptions.Params["key"] = *i.Key
	}
	if i.Provider != nil && *i.Provider != "" {
		requestOptions.Params["provider"] = *i.Provider
	}
	if i.Model != nil && *i.Model != "" {
		requestOptions.Params["model"] = *i.Model
	}
	if i.From != nil {
		requestOptions.Params["from"] = i.From.Format(time.RFC3339)
	}
	if i.To != nil {
		requestOptions.Params["to"] = i.To.Format(time.RFC3339)
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

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "sessions")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *Sessions
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return s, nil
}
