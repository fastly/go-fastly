package usagemetrics

import (
	"context"
	"fmt"
	"io"

	"github.com/fastly/go-fastly/v17/fastly"
)

// Export returns usage metrics for AI Runtime Control requests as raw CSV
// bytes. It accepts the same filter parameters as List, except Sort, which
// isn't supported here. In the CSV output, the column named `time`
// corresponds to the Date field returned by List.
func Export(ctx context.Context, c *fastly.Client, i *ListInput) ([]byte, error) {
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
