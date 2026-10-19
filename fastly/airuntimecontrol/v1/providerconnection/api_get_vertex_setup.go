package providerconnection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetVertexSetup returns the OIDC issuer URL, audience, and customer ID needed to
// configure GCP workload identity federation before creating a `gcp-iam`
// (Vertex AI) provider connection.
func GetVertexSetup(ctx context.Context, c *fastly.Client) (*VertexSetup, error) {
	path := fastly.ToSafeURL("ai-runtime-control", "v1", "providers", "vertex", "setup")

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *VertexSetup
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return s, nil
}
