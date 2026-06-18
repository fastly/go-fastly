package dnszones

import (
	"context"

	"github.com/fastly/go-fastly/v15/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// ZoneID is the Zone Identifier (UUID) (required).
	ZoneID *string `json:"-"`
}

// Delete deletes a specified DNS Zone.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.ZoneID == nil {
		return fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("dns", "v1", "zones", *i.ZoneID)

	resp, err := c.Delete(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
