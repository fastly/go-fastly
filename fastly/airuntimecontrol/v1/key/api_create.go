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
	// Name is a human-readable name for the virtual key (required).
	Name *string `json:"name"`
	// Model is the AI model identifier (required).
	Model *string `json:"model"`
	// Provider is the AI model provider name (required).
	Provider *string `json:"provider"`
	// UserID is the ID of the user creating the key (required).
	UserID *string `json:"user_id"`
	// ExpiresAt is the expiration timestamp of the key.
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// Create creates a new virtual key for AI service authentication.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*VirtualKeyWithToken, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if i.Model == nil {
		return nil, fastly.ErrMissingModel
	}
	if i.Provider == nil {
		return nil, fastly.ErrMissingProvider
	}
	if i.UserID == nil {
		return nil, fastly.ErrMissingUserID
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
