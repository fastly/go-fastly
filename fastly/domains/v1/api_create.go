package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v11/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Description is the description for the domain.
	Description *string `json:"description"`
	// FQDN is the fully-qualified domain name of the domain (required).
	FQDN *string `json:"fqdn"`
	// ServiceID is the service_id associated with the domain or nil if there
	// is no association (optional)
	ServiceID *string `json:"service_id"`
}

// Create creates a new domain.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Data, error) {
	resp, err := c.PostJSON(ctx, "/domains/v1", i, fastly.CreateRequestOptions())
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
