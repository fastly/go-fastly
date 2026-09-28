package versions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
	"github.com/fastly/go-fastly/v17/fastly/domainmanagement/v1/routingconfigs"
)

// ActivateInput specifies the information needed for the Activate() function
// to perform the operation.
type ActivateInput struct {
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string
	// VersionID is the identifier of the version to reactivate (required).
	VersionID *string
}

// Activate reactivates a previous (inactive) version of the specified routing
// config, making it the active version.
func Activate(ctx context.Context, c *fastly.Client, i *ActivateInput) (*routingconfigs.Data, error) {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return nil, fastly.ErrMissingRoutingConfigID
	}
	if i.VersionID == nil || *i.VersionID == "" {
		return nil, fastly.ErrMissingVersionID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "versions", *i.VersionID, "activate")

	resp, err := c.Post(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var d *routingconfigs.Data
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return d, nil
}
