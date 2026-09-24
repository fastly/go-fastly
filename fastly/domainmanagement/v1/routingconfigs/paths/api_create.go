package paths

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Path is the URL path pattern (max 2048 characters, starts with "/")
	// (required).
	Path *string `json:"path"`
	// RoutingConfigID is the routing config identifier (required).
	RoutingConfigID *string `json:"-"`
}

// Create creates a new path within the specified routing config.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Data, error) {
	if i.RoutingConfigID == nil || *i.RoutingConfigID == "" {
		return nil, fastly.ErrMissingRoutingConfigID
	}
	if i.Path == nil {
		return nil, fastly.ErrMissingPath
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs", *i.RoutingConfigID, "paths")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
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
