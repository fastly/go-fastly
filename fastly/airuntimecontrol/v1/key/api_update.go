package key

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// KeyID is the ID of the virtual key (required).
	KeyID *string `json:"-"`
	// Name is the display name of the virtual key.
	Name *string `json:"name,omitempty"`
	// Model is an updated AI model identifier.
	Model *string `json:"model,omitempty"`
	// Provider is an updated AI model provider.
	Provider *string `json:"provider,omitempty"`
	// ExpiresAt is the date and time the key expires, in RFC 3339 format. `null`
	// if the key doesn't expire.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// Update updates the attributes of a virtual key. Only the fields you include
// are changed.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*VirtualKey, error) {
	if i.KeyID == nil {
		return nil, fastly.ErrMissingKeyID
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "keys", *i.KeyID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vk *VirtualKey
	if err := json.NewDecoder(resp.Body).Decode(&vk); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return vk, nil
}
