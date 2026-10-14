package providerconnection

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v17/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// ID is the ID of the provider connection (required).
	ID *string
}

// Delete deletes a provider connection. Returns a 409 response if the
// connection is still referenced by a virtual key's failover policy.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.ID == nil {
		return fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "provider-connections", *i.ID)

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
