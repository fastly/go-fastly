package providerconnection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Name is the display name of the provider connection (required).
	Name *string `json:"name"`
	// Models is the AI model identifiers this connection is allowed to use
	// (required).
	Models []string `json:"models"`
	// BaseURL is the base URL of the provider's API (required). Present only for
	// `api-key` connections.
	BaseURL *string `json:"base_url"`
	// APIKey is the provider's secret API key (required).
	APIKey *string `json:"api_key"`
}

// Create creates a provider connection. Accounts are limited to 25 provider
// connections. Creating a 26th connection returns a 409 response.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*ProviderConnection, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if len(i.Models) == 0 {
		return nil, fastly.ErrMissingModels
	}
	if i.BaseURL == nil {
		return nil, fastly.ErrMissingBaseURL
	}
	if i.APIKey == nil {
		return nil, fastly.ErrMissingAPIKey
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "provider-connections")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
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
