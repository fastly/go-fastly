package tsigkeys

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// UpdateInput specifies the information needed to update a TSIG key.
type UpdateInput struct {
	// TSIGKeyID is the TSIG Key Identifier (UUID) (required).
	TSIGKeyID *string `json:"-"`
	// Name is the name of the TSIG key.
	Name *string `json:"name,omitempty"`
	// Description is a freeform descriptive note.
	Description *fastly.Nullable[string] `json:"description,omitempty"`
	// Algorithm is the algorithm of the TSIG key.
	Algorithm *string `json:"algorithm,omitempty"`
	// Secret is the Base64 encoded secret key.
	Secret *string `json:"secret,omitempty"`
}

// Update updates an existing TSIG key.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*TSIGKey, error) {
	if i.TSIGKeyID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("dns", "v1", "tsig-keys", *i.TSIGKeyID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tsigKey *TSIGKey
	if err := json.NewDecoder(resp.Body).Decode(&tsigKey); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return tsigKey, nil
}
