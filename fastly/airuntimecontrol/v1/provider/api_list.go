package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// List returns the list of AI providers supported by ARC, with each provider's
// available models nested within.
func List(ctx context.Context, c *fastly.Client) (*Providers, error) {
	path := fastly.ToSafeURL("ai-runtime-control", "v1", "providers")

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var p *Providers
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return p, nil
}
