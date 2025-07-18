package redactions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// GetInput specifies the information needed for the Get() function to
// perform the operation.
type GetInput struct {
	// RedactionID is the redaction identifier (required).
	RedactionID *string
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string
}

// Get retrieves the specified redaction.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*Redaction, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.RedactionID == nil {
		return nil, fastly.ErrMissingRedactionID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "redactions", *i.RedactionID)

	resp, err := c.Get(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var redaction *Redaction
	if err := json.NewDecoder(resp.Body).Decode(&redaction); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return redaction, nil
}
