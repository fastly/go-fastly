package accesskeys

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/fastly/go-fastly/v9/fastly"
)

// Permissions can be one of these values only
const (
	READ_WRITE_ADMIN   = "read-write-admin"
	READ_ONLY_ADMIN    = "read-only-admin"
	READ_WRITE_OBJECTS = "read-write-objects"
	READ_ONLY_OBJECTS  = "read-only-objects"
)

var PERMISSONS = []string{READ_WRITE_ADMIN, READ_ONLY_ADMIN, READ_WRITE_OBJECTS, READ_ONLY_OBJECTS}

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Description is a description of the access key (required).
	Description *string `json:"description"`
	// Permission is the permissions the access key will have (required).
	Permission *string `json:"permission"`
	// Buckets are the buckets the access key will have (optional).
	Buckets *[]string `json:"buckets"`
}

// Create creates a new Object Storage Access Key.
func Create(c *fastly.Client, i *CreateInput) (*AccessKey, error) {
	if i.Description == nil {
		return nil, fastly.ErrMissingDescription
	}

	if i.Permission == nil {
		return nil, fastly.ErrMissingPermission
	}

	// Check if the provided permission is in the set of valid permissions
	if !slices.Contains(PERMISSONS, *i.Permission) {
		return nil, fastly.ErrInvalidPermission
	}

	resp, err := c.PostJSON("/resources/object-storage/access-keys", i, nil)
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
