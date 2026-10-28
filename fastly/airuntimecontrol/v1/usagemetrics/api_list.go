package usagemetrics

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/fastly/go-fastly/v18/fastly"
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
	// `sort=-date`). One of `date`, `-date`.
	Sort *string
}

// requestOptions builds the query parameters shared with Export.
func (i *ListInput) requestOptions() fastly.RequestOptions {
	return (&ExportInput{
		Key:      i.Key,
		Provider: i.Provider,
		Model:    i.Model,
		From:     i.From,
		To:       i.To,
		Sort:     i.Sort,
	}).requestOptions()
}

// List returns usage metrics for AI Runtime Control requests, including request
// counts, session counts, input token counts, and output token counts. It
// automatically paginates through all pages returned by the API.
func List(ctx context.Context, c *fastly.Client, i *ListInput) ([]UsageMetric, error) {
	var (
		out    []UsageMetric
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

// listPage retrieves a single page of usage metrics.
func listPage(ctx context.Context, c *fastly.Client, i *ListInput, cursor *string) (*UsageMetrics, error) {
	requestOptions := i.requestOptions()
	if i.Limit != nil {
		requestOptions.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if cursor != nil && *cursor != "" {
		requestOptions.Params["cursor"] = *cursor
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "usage-metrics")

	resp, err := c.Get(ctx, path, requestOptions)
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
