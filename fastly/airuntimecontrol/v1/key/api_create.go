package key

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fastly/go-fastly/v17/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Name is the display name of the virtual key (required).
	Name *string `json:"name"`
	// Model is the AI model identifier.
	Model *string `json:"model,omitempty"`
	// Provider is the AI model provider name (required).
	Provider *string `json:"provider"`
	// ExpiresAt is the date and time the key expires, in RFC 3339 format. `null`
	// if the key doesn't expire.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// SecurityEnabled is whether AI Firewall security is enabled for requests made
	// with this key.
	SecurityEnabled *bool `json:"security_enabled,omitempty"`
	// SecurityAction is the action AI Firewall takes when it detects a violation
	// on a request made with this key. One of `log`, `block`, `replace`.
	SecurityAction *string `json:"security_action,omitempty"`
	// RPMLimit is the requests-per-minute limit for this key. Omit for no limit.
	RPMLimit *int `json:"rpm_limit,omitempty"`
	// TPMLimit is the tokens-per-minute limit for this key. Omit for no limit.
	TPMLimit *int `json:"tpm_limit,omitempty"`
}

// Create creates a new virtual key for authenticating requests through the AI
// Runtime Control gateway.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*VirtualKeyWithToken, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if i.Provider == nil {
		return nil, fastly.ErrMissingProvider
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "keys")

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
