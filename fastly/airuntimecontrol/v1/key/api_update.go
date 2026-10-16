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
	// SecurityEnabled is whether AI Firewall security is enabled for requests made
	// with this key.
	SecurityEnabled *bool `json:"security_enabled,omitempty"`
	// SecurityAction is the action AI Firewall takes when it detects a violation
	// on a request made with this key. One of `log`, `block`, `replace`.
	SecurityAction *string `json:"security_action,omitempty"`
	// RPMLimit is the requests-per-minute limit for this key. Use
	// fastly.NullValue[int]() to remove the limit.
	RPMLimit *fastly.Nullable[int] `json:"rpm_limit,omitempty"`
	// TPMLimit is the tokens-per-minute limit for this key. Use
	// fastly.NullValue[int]() to remove the limit.
	TPMLimit *fastly.Nullable[int] `json:"tpm_limit,omitempty"`
	// SpendLimitEnabled is whether a spend limit is enabled for this key. This
	// field cannot be set when the key is created; it can only be changed with
	// Update.
	SpendLimitEnabled *bool `json:"spend_limit_enabled,omitempty"`
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
