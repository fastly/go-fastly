package timeseries

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v10/fastly"
)

// GetInput specifies the information needed for the Get() function to perform
// the operation.
type GetInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// From is a time range and is the older of the two dates in RFC 3339 format (required).
	From *string
	// Granularity is the sample size in seconds (optional).
	Granularity *int
	// Metrics are comma separated list of metrics to be included in the timeseries (required).
	Metrics *string
	// To is a time range and is the older of the two dates in RFC 3339 format (optional).
	To *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Get retrieves the specified workspace.
func Get(c *fastly.Client, i *GetInput) (*TimeSeries, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	if i.Metrics == nil {
		return nil, fastly.ErrMissingMetrics
	}

	if i.From == nil {
		return nil, fastly.ErrMissingFrom
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "timeseries")

	requestOptions := fastly.CreateRequestOptions(i.Context)
	if i.From != nil {
		requestOptions.Params["from"] = *i.From
	}
	if i.To != nil {
		requestOptions.Params["to"] = *i.To
	}
	if i.Metrics != nil {
		requestOptions.Params["metrics"] = *i.Metrics
	}
	if i.Granularity != nil {
		requestOptions.Params["granularity"] = strconv.Itoa(*i.Granularity)
	}

	resp, err := c.GetJSON(path, requestOptions)
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
