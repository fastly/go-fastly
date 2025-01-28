package computeacls

import (
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v9/fastly"
)

// ListACLs retrieves all compute ACLs.
func ListACLs(c *fastly.Client) (*ComputeACLs, error) {
	resp, err := c.Get("/resources/acls", nil)
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
