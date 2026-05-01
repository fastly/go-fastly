package dnszones

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// GetInput specifies the information needed to get a zone.
type GetInput struct {
	// ZoneID is the Zone Identifier (UUID) (required).
	ZoneID *string `json:"-"`
}

// Get retrieves a specified DNS Zone.
func Get(ctx context.Context, c *fastly.Client, i *GetInput) (*Zone, error) {
	if i.ZoneID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("dns", "v1", "zones", *i.ZoneID)

	resp, err := c.GetJSON(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var zone *Zone
	if err := json.NewDecoder(resp.Body).Decode(&zone); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return zone, nil
}
