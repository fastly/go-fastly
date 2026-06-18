package timeseries

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v15/fastly"
)

// GetInput specifies the information needed for the Get() function to
// perform the operation.
type GetInput struct {
	// End is the end of a date-time range, expressed in RFC 3339 format.
	End *string
	// Granularity is the level of detail of the sample size in seconds.
	Granularity *int
	// Metrics is a comma-separated list of metrics to be included in the
	// timeseries. Metrics can be XSS, SQLI, HTTP404, requests_total,
	// requests_attack, requests_total_blocked, or any custom metric (required).
	Metrics *string
	// Start is the start of a date-time range, expressed in RFC 3339 format (required).
	Start *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Get retrieves the specified timeseries.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*TimeSeries, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	if i.Metrics == nil {
		return nil, fastly.ErrMissingMetrics
	}

	if i.Start == nil {
		return nil, fastly.ErrMissingStart
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "timeseries")

	requestOptions := fastly.CreateRequestOptions()
	if i.Start != nil {
		requestOptions.Params["start"] = *i.Start
	}
	if i.End != nil {
		requestOptions.Params["end"] = *i.End
	}
	if i.Metrics != nil {
		requestOptions.Params["metrics"] = *i.Metrics
	}
	if i.Granularity != nil {
		requestOptions.Params["granularity"] = strconv.Itoa(*i.Granularity)
	}

	resp, err := c.GetJSON(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ts *TimeSeries
	if err := json.NewDecoder(resp.Body).Decode(&ts); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return ts, nil
}
