package tsigkeys

import (
	"context"

	"github.com/fastly/go-fastly/v14/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// TSIGKeyID is the TSIG Key Identifier (UUID) (required).
	TSIGKeyID *string `json:"-"`
}

// Delete deletes a specified TSIG key.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.TSIGKeyID == nil {
		return fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("dns", "v1", "tsig-keys", *i.TSIGKeyID)

	resp, err := c.Delete(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
