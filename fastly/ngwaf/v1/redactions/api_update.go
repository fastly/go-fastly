package redactions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Field is the name of the field to redact (required). Will be converted to lowercase.
	Field *string `json:"field"`
	// RedactionID is the id of the redaction that's being updated (required).
	RedactionID *string `json:"-"`
	// Type is the type of field being redacted. Must be one of `request_parameter`, `request_header`, or `response_header`.
	Type *string `json:"type"`
	// WorkspaceID is the ID of the workspace that the redaction belongs to.
	WorkspaceID *string `json:"-"`
}

// Update updates the specified workspace.
func Update(c *fastly.Client, i *UpdateInput) (*Redaction, error) {
	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}
	if i.RedactionID == nil {
		return nil, fastly.ErrMissingRedactionID
	}
	if i.Field == nil {
		return nil, fastly.ErrMissingField
	}
	if i.Type == nil {
		return nil, fastly.ErrMissingType
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "redactions", *i.RedactionID)

	resp, err := c.PatchJSON(path, i, fastly.CreateRequestOptions(i.Context))
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
