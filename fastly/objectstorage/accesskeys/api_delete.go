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
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
}

// DeleteAccessKey deletes an access key.
func Delete(c *fastly.Client, i *DeleteInput) error {
	if i.AccessKeyID == nil {
		return fastly.ErrMissingAccessKeyID
	}

	path := fastly.ToSafeURL("resources", "object-storage", "access-keys", *i.AccessKeyID)

	resp, err := c.Delete(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fastly.NewHTTPError(resp)
	}
	return nil
}
