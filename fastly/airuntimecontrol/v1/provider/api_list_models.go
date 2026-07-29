package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// ListModelsInput specifies the information needed for the ListModels()
// function to perform the operation.
type ListModelsInput struct {
	// ProviderID is the ID identifying the provider (required).
	ProviderID *string
}

// ListModels returns the list of models available for a specific provider.
func ListModels(ctx context.Context, c *fastly.Client, i *ListModelsInput) (*Models, error) {
	if i.ProviderID == nil {
		return nil, fastly.ErrMissingProviderID
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "providers", *i.ProviderID, "models")

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var m *Models
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return m, nil
}
