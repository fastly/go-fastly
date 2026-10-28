package slack

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fastly/go-fastly/v18/fastly"
	ngwafv1 "github.com/fastly/go-fastly/v18/fastly/ngwaf/v1"
)

// ListInput specifies the information needed for the List() function to perform
// the operation.
type ListInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// Limit how many results are returned. Defaults to the API maximum
	// of 1000.
	Limit *int
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// List retrieves a list of slack workspace alerts.
func List(ctx context.Context, c *fastly.Client, i *ListInput) (*Alerts, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	requestOptions := fastly.CreateRequestOptions()
	limit := ngwafv1.MaxLimit
	if i.Limit != nil {
		limit = *i.Limit
	}
	requestOptions.Params["limit"] = strconv.Itoa(limit)

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "alerts")

	resp, err := c.Get(ctx, path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var was *Alerts
	if err := json.NewDecoder(resp.Body).Decode(&was); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	// Parse the alerts to only include the desired type of integration.
	var parsedAlerts []Alert
	for _, alert := range was.Data {
		if alert.Type == IntegrationType {
			parsedAlerts = append(parsedAlerts, alert)
		}
	}
	was.Data = parsedAlerts
	was.Meta.Total = len(parsedAlerts)

	return was, nil
}
