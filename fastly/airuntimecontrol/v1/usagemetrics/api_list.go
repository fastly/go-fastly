package usagemetrics

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListInput specifies the information needed for the List() and Export()
// functions to perform the operation.
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
	// Sort is the sort field. Prefix with `-` for descending order (e.g.,
	// `sort=-date`). One of `date`, `-date`.
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
	if i.Sort != nil && *i.Sort != "" {
		requestOptions.Params["sort"] = *i.Sort
	}
	return requestOptions
}

// List returns usage metrics for AI Runtime Control requests, including request
// counts, session counts, input token counts, and output token counts.
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
