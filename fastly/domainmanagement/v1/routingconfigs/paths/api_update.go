package paths

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// Path is the URL path pattern (max 2048 characters, starts with "/").
	Path *string `json:"path,omitempty"`
	// PathID is the path identifier (required).
	PathID *string `json:"-"`
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string `json:"-"`
}

// Update updates the specified path.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Data, error) {
	if i.RoutingConfigID == nil {
		return nil, fastly.ErrMissingRoutingConfigID
	}
	if i.PathID == nil {
		return nil, fastly.ErrMissingPathID
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "paths", *i.PathID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
