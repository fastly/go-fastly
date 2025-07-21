package signals

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/common"
)

// CreateInput specifies the information needed for the Create()
// function to perform the operation.
type CreateInput struct {
	// Description is a description of the signal.
	Description *string `json:"description,omitempty"`
	// Name is the name of the signal. Must be between 3 and 25
	// characters. Letters, numbers, hyphens, and spaces are
	// accepted. Special characters and periods are not accepted (required).
	Name *string `json:"name"`
	// Scope defines where the signal is located, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *common.Scope `json:"-"`
}

// Create creates a new signal in the given workspace.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Signal, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if i.Scope == nil {
		return nil, fastly.ErrMissingScope
	}

	path, err := common.BuildPath(i.Scope, "signals", "")
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var signal *Signal
	if err := json.NewDecoder(resp.Body).Decode(&signal); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return signal, nil
}
