package fastly

import (
	"fmt"
	"strconv"
	"time"
)

// Resource represents a response from the Fastly API.
type Resource struct {
	// CreatedAt is the date and time in ISO 8601 format.
	CreatedAt *time.Time `mapstructure:"created_at" json:"created_at"`
	// CreatedAt is the date and time in ISO 8601 format.
	DeletedAt *time.Time `mapstructure:"deleted_at" json:"deleted_at"`
	// HREF is the path to the resource.
	HREF *string `mapstructure:"href" json:"href"`
	// LinkID is an alphanumeric string identifying the resource link.
	LinkID *string `mapstructure:"id" json:"id"`
	// Name is the name of the resource being linked to.
	Name *string `mapstructure:"name" json:"name"`
	// ResourceID is the ID of the linked resource.
	ResourceID *string `mapstructure:"resource_id" json:"resource_id"`
	// ResourceType is the type of the linked resource.
	ResourceType *string `mapstructure:"resource_type" json:"resource_type"`
	// ServiceID is an alphanumeric string identifying the service.
	ServiceID *string `mapstructure:"service_id" json:"service_id"`
	// ServiceVersion is an integer identifying a service version.
	ServiceVersion *int `mapstructure:"version" json:"version"`
	// UpdatedAt is the date and time in ISO 8601 format.
	UpdatedAt *time.Time `mapstructure:"updated_at" json:"updated_at"`
}

// ListResourcesInput is used as input to the ListResources function.
type ListResourcesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListResources retrieves all resources.
func (c *Client) ListResources(i *ListResourcesInput) ([]*Resource, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "resource")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rs []*Resource
	if err := decodeBodyMap(resp.Body, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

// CreateResourceInput is used as input to the CreateResource function.
type CreateResourceInput struct {
	// Name is the name of the resource being linked to (e.g. a kv store).
	//
	// NOTE: This doesn't have to match the actual resource name, i.e. the name
	// of the KV Store. Rather, this is an opportunity for you to use an
	// 'alias' for your KV Store. So your service will now refer to the
	// KV Store using this name.
	Name *string `url:"name,omitempty"`
	// ResourceID is the ID of the linked resource.
	ResourceID *string `url:"resource_id,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// CreateResource creates a new resource.
func (c *Client) CreateResource(i *CreateResourceInput) (*Resource, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "resource")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r *Resource
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// GetResourceInput is used as input to the GetResource function.
type GetResourceInput struct {
	// ResourceID is an alphanumeric string identifying the resource link (required).
	ResourceID string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetResource retrieves the specified resource.
func (c *Client) GetResource(i *GetResourceInput) (*Resource, error) {
	if i.ResourceID == "" {
		return nil, ErrMissingResourceID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "resource", i.ResourceID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r *Resource
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// UpdateResourceInput is used as input to the UpdateResource function.
type UpdateResourceInput struct {
	// Name is the name of the resource being linked to (e.g. a kv store).
	Name *string `url:"name,omitempty"`
	// ResourceID is an alphanumeric string identifying the resource link (required).
	ResourceID string `url:"-"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// UpdateResource updates the specified resource.
func (c *Client) UpdateResource(i *UpdateResourceInput) (*Resource, error) {
	if i.ResourceID == "" {
		return nil, ErrMissingResourceID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "resource", i.ResourceID)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Resource
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteResourceInput is the input parameter to DeleteResource.
type DeleteResourceInput struct {
	// ResourceID is an alphanumeric string identifying the resource link (required).
	ResourceID string `url:"-"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// DeleteResource deletes the specified resource.
func (c *Client) DeleteResource(i *DeleteResourceInput) error {
	if i.ResourceID == "" {
		return ErrMissingResourceID
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "resource", i.ResourceID)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return fmt.Errorf("not ok")
	}
	return nil
}
