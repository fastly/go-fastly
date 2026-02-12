package signals

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v13/fastly"
	"github.com/fastly/go-fastly/v13/fastly/ngwaf/v1/scope"
)

// GetInput specifies the information needed for the Get() function to
// perform the operation.
type GetInput struct {
	// Scope defines where the signal is applied, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *scope.Scope
	// SignalID is the signal identifier (required).
	SignalID *string
}

// Get retrieves the specified signal for the given workspace.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*Signal, error) {
	if i.Scope == nil {
		return nil, fastly.ErrMissingScope
	}
	if i.SignalID == nil {
		return nil, fastly.ErrMissingSignalID
	}

	path, err := scope.BuildPath(i.Scope, "signals", *i.SignalID)
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
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
