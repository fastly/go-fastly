package usagemetrics

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/fastly/go-fastly/v18/fastly"
)

// ExportInput specifies the information needed for the Export() function to
// perform the operation.
type ExportInput struct {
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

// requestOptions builds the query parameters for the usage metrics list and
// export endpoints.
func (i *ExportInput) requestOptions() fastly.RequestOptions {
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

// Export returns usage metrics for AI Runtime Control requests as raw CSV
// bytes. Unlike List, it does not paginate. In the CSV output, the column
// named `time` corresponds to the Date field returned by List.
func Export(ctx context.Context, c *fastly.Client, i *ExportInput) ([]byte, error) {
	requestOptions := i.requestOptions()
	requestOptions.Headers["Accept"] = "text/csv"

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "usage-metrics", "export")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
