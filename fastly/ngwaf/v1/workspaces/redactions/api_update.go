package redactions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v13/fastly"
)

// UpdateInput specifies the information needed for the Update()
// function to perform the operation.
type UpdateInput struct {
	// Field is the name of the field to redact. Will be converted
	// to lowercase.
	Field *string `json:"field,omitempty"`
	// RedactionID is the id of the redaction that's being updated
	// (required).
	RedactionID *string `json:"-"`
	// Type is the type of field being redacted. Must be one of
	// `request_parameter`, `request_header`, or
	// `response_header`.
	Type *string `json:"type,omitempty"`
	// WorkspaceID is the ID of the workspace that the redaction
	// belongs to (required).
	WorkspaceID *string `json:"-"`
}

// Update updates the specified redaction.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*Redaction, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.RedactionID == nil {
		return nil, fastly.ErrMissingRedactionID
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "redactions", *i.RedactionID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
