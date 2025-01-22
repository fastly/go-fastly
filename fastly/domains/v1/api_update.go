package v1

import (
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v9/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// DomainID is the domain identifier (required).
	DomainID *string `json:"-"`
	// ServiceID is the service_id associated with the domain or nil if there
	// is no association (optional)
	ServiceID *string `json:"service_id"`
}

// Update updates the specified domain.
func Update(c *fastly.Client, i *UpdateInput) (*Data, error) {
	if i.DomainID == nil {
		return nil, fastly.ErrMissingDomainID
	}

	path := fastly.ToSafeURL("domains", "v1", *i.DomainID)

	resp, err := c.PatchJSON(path, i, nil)
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
