package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// GetInput specifies the information needed for the Get() function to perform
// the operation.
type GetInput struct {
	// DomainID is the domain identifier (required).
	DomainID *string
}

// Get retrieves a specified domain.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*Data, error) {
	if i.DomainID == nil {
		return nil, fastly.ErrMissingDomainID
	}

	path := fastly.ToSafeURL("domains", "v1", *i.DomainID)

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
