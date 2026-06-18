package timeseries

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v15/fastly"
)

// ListInput specifies the information needed for the List() function to
// perform the operation.
type ListInput struct {
	// Dimensions is a comma-separated list of grouping dimensions to be
	// included in the timeseries. Allowed values are workspaces and time.
	// Default is time.
	Dimensions *string
	// From is the start of a date-time range, expressed in RFC 3339 format (required).
	From *string
	// Granularity is the level of detail of the sample size in seconds.
	Granularity *int
	// Metrics is a comma-separated list of metrics to be included in the
	// timeseries. Metrics can be XSS, SQLI, HTTP404, requests_total,
	// requests_attack, requests_total_blocked, or any custom metric (required).
	Metrics *string
	// To is the end of a date-time range, expressed in RFC 3339 format.
	To *string
}

// List retrieves timeseries metrics for Next-Gen WAF.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Timeseries, error) {
	if i.Metrics == nil || *i.Metrics == "" {
		return nil, fastly.ErrMissingMetrics
	}

	if i.From == nil || *i.From == "" {
		return nil, fastly.ErrMissingFrom
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "timeseries")

	requestOptions := fastly.CreateRequestOptions()
	requestOptions.Params["metrics"] = *i.Metrics
	requestOptions.Params["from"] = *i.From
	if i.To != nil {
		requestOptions.Params["to"] = *i.To
	}
	if i.Dimensions != nil {
		requestOptions.Params["dimensions"] = *i.Dimensions
	}
	if i.Granularity != nil {
		requestOptions.Params["granularity"] = strconv.Itoa(*i.Granularity)
	}

	resp, err := c.GetJSON(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ts *Timeseries
	if err := json.NewDecoder(resp.Body).Decode(&ts); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return ts, nil
}
