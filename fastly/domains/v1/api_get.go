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
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
	// DomainID is the domain identifier (required).
	DomainID *string
}

// Get retrieves a specified domain.
func Get(c *fastly.Client, i *GetInput) (*Data, error) {
	if i.DomainID == nil {
		return nil, fastly.ErrMissingDomainID
	}

	path := fastly.ToSafeURL("domains", "v1", *i.DomainID)

	resp, err := c.Get(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer fastly.CheckCloseForErr(resp.Body.Close)

	var d *Data
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return d, nil
}
