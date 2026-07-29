package key

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// GetInput specifies the information needed for the Get() function to perform
// the operation.
type GetInput struct {
	// KeyID is the ID identifying the virtual key (required).
	KeyID *string
}

// Get retrieves information on a specific virtual key.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*VirtualKeyListItem, error) {
	if i.KeyID == nil {
		return nil, fastly.ErrMissingKeyID
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "keys", *i.KeyID)

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vk *VirtualKeyListItem
	if err := json.NewDecoder(resp.Body).Decode(&vk); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return vk, nil
}
