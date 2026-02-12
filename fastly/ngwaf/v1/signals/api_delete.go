package signals

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fastly/go-fastly/v13/fastly"
	"github.com/fastly/go-fastly/v13/fastly/ngwaf/v1/scope"
)

// DeleteInput specifies the information needed for the Delete()
// function to perform the operation.
type DeleteInput struct {
	// Scope defines where the signal is applied, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *scope.Scope
	// SignalID is the signal identifier (required).
	SignalID *string
}

// Delete deletes the specified signal.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.SignalID == nil {
		return fastly.ErrMissingSignalID
	}
	if i.Scope == nil {
		return fastly.ErrMissingScope
	}

	path, err := scope.BuildPath(i.Scope, "signals", *i.SignalID)
	if err != nil {
		return fmt.Errorf("failed to build API path: %w", err)
	}

	resp, err := c.Delete(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
