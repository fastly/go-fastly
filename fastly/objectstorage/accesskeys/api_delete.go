package accesskeys

import (
	"context"
	"net/http"

	"github.com/fastly/go-fastly/v10/fastly"
)

// DeleteInput specifies the information needed for the Delete() function to
// perform the operation.
type DeleteInput struct {
	// AccessKeyID is an AccessKey Identifier (required).
	AccessKeyID *string
}

// DeleteAccessKey deletes an access key.
func Delete(ctx context.Context, c *fastly.Client, i *DeleteInput) error {
	if i.AccessKeyID == nil {
		return fastly.ErrMissingAccessKeyID
	}

	path := fastly.ToSafeURL("resources", "object-storage", "access-keys", *i.AccessKeyID)

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
