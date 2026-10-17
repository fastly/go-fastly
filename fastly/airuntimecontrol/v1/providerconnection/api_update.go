package providerconnection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// ID is the ID of the provider connection (required).
	ID *string `json:"-"`
	// Models is the AI model identifiers this connection is allowed to use.
	Models []string `json:"models,omitempty"`
	// BaseURL is the base URL of the provider's API. Present only for `api-key`
	// connections.
	BaseURL *string `json:"base_url,omitempty"`
	// APIKey is the provider's secret API key.
	APIKey *string `json:"api_key,omitempty"`
}

// Update updates the attributes of a provider connection. Only the fields you
// include are changed.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*ProviderConnection, error) {
	if i.ID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "provider-connections", *i.ID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
