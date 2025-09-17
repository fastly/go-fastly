package lists

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fastly/go-fastly/v12/fastly"
	"github.com/fastly/go-fastly/v12/fastly/ngwaf/v1/scope"
)

// DeleteInput specifies the information needed for the Delete()
// function to perform the operation.
type DeleteInput struct {
	// ListID is the id of the list to be deleted (required).
	ListID *string
	// Scope defines where the list is applied, including its type (e.g.,
	// "workspace" or "account") and the specific IDs it applies to (required).
	Scope *scope.Scope
}

// Delete deletes the specified list.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.ListID == nil {
		return fastly.ErrMissingListID
	}
	if i.Scope == nil {
		return fastly.ErrMissingScope
	}

	path, err := scope.BuildPath(i.Scope, "lists", *i.ListID)
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
