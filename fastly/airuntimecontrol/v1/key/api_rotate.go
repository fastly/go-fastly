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
	// KeyID is the ID of the virtual key (required).
	KeyID *string `json:"-"`
	// ExpiresAt is the date and time the rotated key expires, in RFC 3339 format.
	// If omitted, the rotated key does not expire, even if the key previously had
	// an expiration.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// Rotate issues a new access token for a virtual key, invalidating the previous
// one.
func Rotate(ctx context.Context, c *fastly.Client, i *RotateInput) (*VirtualKeyWithToken, error) {
	if i.KeyID == nil {
		return nil, fastly.ErrMissingKeyID
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
