package tsigkeys

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// GetInput specifies the information needed to get a TSIG key.
type GetInput struct {
	// TSIGKeyID is the TSIG Key Identifier (UUID) (required).
	TSIGKeyID *string `json:"-"`
}

// Get retrieves a specified TSIG key.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*TSIGKey, error) {
	if i.TSIGKeyID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("dns", "v1", "tsig-keys", *i.TSIGKeyID)

	resp, err := c.GetJSON(ctx, path, fastly.CreateRequestOptions())
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
