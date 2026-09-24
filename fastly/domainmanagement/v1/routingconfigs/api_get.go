package routingconfigs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetInput specifies the information needed for the Get() function to perform
// the operation.
type GetInput struct {
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
}

// Get retrieves a specified routing config.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*Data, error) {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return nil, fastly.ErrMissingRoutingConfigID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID)

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
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
