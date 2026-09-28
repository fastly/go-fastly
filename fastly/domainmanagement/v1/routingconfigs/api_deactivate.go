package routingconfigs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// DeactivateInput specifies the information needed for the Deactivate()
// function to perform the operation.
type DeactivateInput struct {
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
}

// Deactivate deactivates the specified routing config's active version,
// leaving it with no active version.
func Deactivate(ctx context.Context, c *fastly.Client, i *DeactivateInput) (*Data, error) {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return nil, fastly.ErrMissingRoutingConfigID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "deactivate")

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
