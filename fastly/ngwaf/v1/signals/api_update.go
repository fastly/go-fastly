package signals

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/common"
)

// UpdateInput specifies the information needed for the Update()
// function to perform the operation.
type UpdateInput struct {
	// Description is the new description for the signal
	// (required).
	Description *string `json:"description"`
	// Scope defines where the signal is located, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *common.Scope
	// SignalID is the id of the signal that's being updated
	// (required).
	SignalID *string `json:"-"`
}

// Update updates the specified signal.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Signal, error) {
	if i.SignalID == nil {
		return nil, fastly.ErrMissingSignalID
	}
	if i.Scope == nil {
		return nil, fastly.ErrMissingScope
	}
	if i.Description == nil {
		return nil, fastly.ErrMissingDescription
	}

	path, err := common.BuildPath(i.Scope, "signals", *i.SignalID)
	if err != nil {
		return nil, fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
