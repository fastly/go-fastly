package routingconfigs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Name is the user-defined name of the routing config (required).
	Name *string `json:"name"`
}

// Create creates a new routing config.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Data, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}

	path := fastly.ToSafeURL("domain-management", "v1", "routing-configs")

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
