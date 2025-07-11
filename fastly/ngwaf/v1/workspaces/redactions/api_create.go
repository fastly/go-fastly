package redactions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Field is the name of the field to redact (required). Will be converted to lowercase.
	Field *string `json:"field"`
	// Type is the type of field being redacted. Must be one of `request_parameter`, `request_header`, or `response_header`.
	Type *string `json:"type"`
	// WorkspaceID is the ID of the workspace that the redaction is being created in.
	WorkspaceID *string `json:"-"`
}

// Create creates a new redaction.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*Redaction, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.Field == nil {
		return nil, fastly.ErrMissingField
	}
	if i.Type == nil {
		return nil, fastly.ErrMissingType
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "redactions")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
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
