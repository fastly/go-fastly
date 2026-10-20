package key

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v18/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// KeyID is the ID of the virtual key (required).
	KeyID *string
}

// Delete deletes a virtual key. This also deletes any security configuration
// associated with the key.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.KeyID == nil {
		return fastly.ErrMissingKeyID
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "keys", *i.KeyID)

	resp, err := c.Delete(ctx, path, fastly.CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fastly.NewHTTPError(resp)
	}

	return nil
}
