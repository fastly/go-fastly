package tsigkeys

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// CreateInput specifies the information needed to create a TSIG key.
type CreateInput struct {
	// Name is the name of the TSIG key (required).
	Name *string `json:"name"`
	// Algorithm is the algorithm of the TSIG key (required).
	Algorithm *string `json:"algorithm"`
	// Secret is the Base64 encoded secret key (required).
	Secret *string `json:"secret"`
	// Description is a freeform descriptive note.
	Description *string `json:"description,omitempty"`
}

// Create creates a new TSIG key.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*TSIGKey, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if i.Algorithm == nil {
		return nil, fastly.ErrMissingAlgorithm
	}
	if i.Secret == nil {
		return nil, fastly.ErrMissingSecret
	}

	resp, err := c.PostJSON(ctx, "/dns/v1/tsig-keys", i, fastly.CreateRequestOptions())
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
