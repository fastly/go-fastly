package routingconfigs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ActivateInput specifies the information needed for the Activate() function
// to perform the operation.
type ActivateInput struct {
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
}

// Activate activates the current draft version of the specified routing
// config, making it the active version.
func Activate(ctx context.Context, c *fastly.Client, i *ActivateInput) (*Data, error) {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return nil, fastly.ErrMissingRoutingConfigID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "activate")

	resp, err := c.Post(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *Data
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return d, nil
}
