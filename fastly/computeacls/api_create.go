package computeacls

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Name is the name of the compute ACL to create (required).
	Name *string `json:"name"`
}

// Create creates a new compute ACL.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*ComputeACL, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}

	resp, err := c.PostJSON(ctx, "/resources/acls", i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var acl *ComputeACL
	if err := json.NewDecoder(resp.Body).Decode(&acl); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return acl, nil
}
