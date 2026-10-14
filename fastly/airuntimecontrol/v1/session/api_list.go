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
	// Key filters results to a specific virtual key ID.
	Key *string
	// Provider filters results by AI provider connection ID.
	Provider *string
	// Model filters results by AI model identifier.
	Model *string
	// From is the start of the time range, in RFC 3339 format.
	From *time.Time
	// To is the end of the time range, in RFC 3339 format. Defaults to `now`.
	To *time.Time
	// Limit limits how many results are returned per page (1-1000).
	Limit *int
	// Sort is the sort field. Prefix with `-` for descending order (e.g.,
	// `sort=-created_at`). One of `created_at`, `-created_at`.
	Sort *string
}

// List returns session logs, including request and response data and metadata,
// for debugging and audit purposes. It automatically paginates through all
// pages returned by the API.
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
