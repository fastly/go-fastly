package providerconnection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetInput specifies the information needed for the Get() function to perform
// the operation.
type GetInput struct {
	// ID is the ID identifying the provider connection (required).
	ID *string
}

// Get retrieves a specific provider connection for a customer.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*ProviderConnection, error) {
	if i.ID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "provider-connections", *i.ID)

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pc *ProviderConnection
	if err := json.NewDecoder(resp.Body).Decode(&pc); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return pc, nil
}
