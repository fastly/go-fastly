package dnszones

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v14/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// ZoneID is the Zone Identifier (UUID) (required).
	ZoneID *string `json:"-"`
}

// Delete deletes a specified DNS Zone.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) (*Zone, error) {

	if i.ZoneID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("dns", "v1", "zones", *i.ZoneID)

	resp, err := c.Delete(ctx, path, fastly.CreateRequestOptions())
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
