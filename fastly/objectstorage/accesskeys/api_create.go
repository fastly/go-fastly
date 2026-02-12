package accesskeys

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/fastly/go-fastly/v13/fastly"
)

// Permissions can be one of these values only.
const (
	ReadWriteAdmin  = "read-write-admin"
	ReadOnlyAdmin   = "read-only-admin"
	ReadWriteObject = "read-write-objects"
	ReadOnlyObjects = "read-only-objects"
)

var PERMISSIONS = []string{ReadWriteAdmin, ReadOnlyAdmin, ReadWriteObject, ReadOnlyObjects}

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Description is a description of the access key (required).
	Description *string `json:"description"`
	// Permission is the permissions the access key will have (required).
	Permission *string `json:"permission"`
	// Buckets are the buckets the access key will have.
	Buckets *[]string `json:"buckets"`
}

// Create creates a new Object Storage Access Key.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*AccessKey, error) {
	if i.Description == nil {
		return nil, fastly.ErrMissingDescription
	}

	if i.Permission == nil {
		return nil, fastly.ErrMissingPermission
	}

	// Check if the provided permission is in the set of valid permissions
	if !slices.Contains(PERMISSIONS, *i.Permission) {
		return nil, fastly.ErrInvalidPermission
	}

	resp, err := c.PostJSON(ctx, "/resources/object-storage/access-keys", i, fastly.CreateRequestOptions())
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
