package key

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fastly/go-fastly/v17/fastly"
)

// RotateInput specifies the information needed for the Rotate() function to
// perform the operation.
type RotateInput struct {
	// KeyID is the ID identifying the virtual key (required).
	KeyID *string `json:"-"`
	// ExpiresAt is the expiration timestamp for the rotated key (required).
	ExpiresAt *time.Time `json:"expires_at"`
}

// Rotate rotates an existing virtual key, generating a new access token.
func Rotate(ctx context.Context, c *fastly.Client, i *RotateInput) (*VirtualKeyWithToken, error) {
	if i.KeyID == nil {
		return nil, fastly.ErrMissingKeyID
	}
	if i.ExpiresAt == nil {
		return nil, fastly.ErrMissingExpiresAt
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "keys", *i.KeyID, "rotate")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vk *VirtualKeyWithToken
	if err := json.NewDecoder(resp.Body).Decode(&vk); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return vk, nil
}
