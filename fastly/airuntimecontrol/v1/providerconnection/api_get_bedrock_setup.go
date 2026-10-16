package providerconnection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetBedrockSetup returns the AI Runtime Control role ARN and external ID needed to
// configure your IAM trust policy before creating an `aws-iam` (Bedrock)
// provider connection.
func GetBedrockSetup(ctx context.Context, c *fastly.Client) (*BedrockSetup, error) {
	path := fastly.ToSafeURL("ai-runtime-control", "v1", "providers", "bedrock", "setup")

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s *BedrockSetup
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return s, nil
}
