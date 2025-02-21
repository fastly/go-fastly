package accesskeys

import (
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v9/fastly"
)

// GetInput specifies the information needed for the Get() function to perform
// the operation.
type GetInput struct {
	// AccessKeyID is an AccessKey Identifier (required).
	AccessKeyID *string
}

// GetAccessKey finds an access key with the given ID if the user has the correct permisssions.
func Get(c *fastly.Client, i *GetInput) (*AccessKey, error) {
	if i.AccessKeyID == nil {
		return nil, fastly.ErrMissingAccessKeyID
	}

	path := fastly.ToSafeURL("resources", "object-storage", "access-keys", *i.AccessKeyID)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var accessKey *AccessKey
	if err := json.NewDecoder(resp.Body).Decode(&accessKey); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return accessKey, nil
}
