package usagemetrics

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List() and Export()
// functions to perform the operation.
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
	// Cursor is the pagination cursor.
	Cursor *string
	// Limit is the maximum number of results per page.
	Limit *int
	// Sort is the sort field. Prefix with "-" for descending order (e.g. "-date").
	Sort *string
}

// requestOptions builds the shared query parameters for the usage metrics
// list and export endpoints.
func (i *ListInput) requestOptions() fastly.RequestOptions {
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
	if i.Cursor != nil && *i.Cursor != "" {
		requestOptions.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Sort != nil && *i.Sort != "" {
		requestOptions.Params["sort"] = *i.Sort
	}
	return requestOptions
}

// List returns usage metrics for AI services, with optional filtering and
// pagination.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*UsageMetrics, error) {
	path := fastly.ToSafeURL("ai-runtime-control", "v1", "usage-metrics")

	resp, err := c.Get(ctx, path, i.requestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var m *UsageMetrics
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return m, nil
}
