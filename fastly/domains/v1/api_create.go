package v1

import (
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v9/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// FQDN is the fully-qualified domain name of the domain (required).
	FQDN *string `json:"fqdn"`
	// ServiceID is the service_id associated with the domain or nil if there
	// is no association (optional)
	ServiceID *string `json:"service_id"`
}

// Create creates a new domain.
func Create(c *fastly.Client, i *CreateInput) (*Data, error) {
	resp, err := c.PostJSON("/domains/v1", i, nil)
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
