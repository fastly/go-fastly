package accesskeys

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// GetInput specifies the information needed for the Get() function to perform
// the operation.
type GetInput struct {
	// AccessKeyID is an AccessKey Identifier (required).
	AccessKeyID *string
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context
}

// GetAccessKey finds an access key with the given ID if the user has the correct permisssions.
func Get(c *fastly.Client, i *GetInput) (*AccessKey, error) {
	if i.AccessKeyID == nil {
		return nil, fastly.ErrMissingAccessKeyID
	}

	path := fastly.ToSafeURL("resources", "object-storage", "access-keys", *i.AccessKeyID)

	resp, err := c.Get(path, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer fastly.CheckCloseForErr(resp.Body.Close)

	var accessKey *AccessKey
	if err := json.NewDecoder(resp.Body).Decode(&accessKey); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return accessKey, nil
}
