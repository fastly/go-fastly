package computeacls

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v12/fastly"
)

// ListACLs retrieves all compute ACLs.
func ListACLs(ctx context.Context, c *fastly.Client) (*ComputeACLs, error) {
	resp, err := c.Get(ctx, "/resources/acls", fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var acls *ComputeACLs

	if err := json.NewDecoder(resp.Body).Decode(&acls); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return acls, nil
}
