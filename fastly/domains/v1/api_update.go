package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v11/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// Description is the description for the domain.
	Description *string `json:"description"`
	// DomainID is the domain identifier (required).
	DomainID *string `json:"-"`
	// ServiceID is the service_id associated with the domain or nil if there
	// is no association (optional)
	ServiceID *string `json:"service_id"`
}

// Update updates the specified domain.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Data, error) {
	if i.DomainID == nil {
		return nil, fastly.ErrMissingDomainID
	}

	path := fastly.ToSafeURL("domains", "v1", *i.DomainID)

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
