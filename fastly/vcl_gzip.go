package fastly

import (
	"context"
	"strconv"
	"time"
)

// Gzip represents an Gzip logging response from the Fastly API.
type Gzip struct {
	CacheCondition *string    `mapstructure:"cache_condition"`
	ContentTypes   *string    `mapstructure:"content_types"`
	CreatedAt      *time.Time `mapstructure:"created_at"`
	DeletedAt      *time.Time `mapstructure:"deleted_at"`
	Extensions     *string    `mapstructure:"extensions"`
	Name           *string    `mapstructure:"name"`
	ServiceID      *string    `mapstructure:"service_id"`
	ServiceVersion *int       `mapstructure:"version"`
	UpdatedAt      *time.Time `mapstructure:"updated_at"`
}

// ListGzipsInput is used as input to the ListGzips function.
type ListGzipsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListGzips retrieves all resources.
func (c *Client) ListGzips(ctx context.Context, i *ListGzipsInput) ([]*Gzip, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "gzip")
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gzips []*Gzip
	if err := DecodeBodyMap(resp.Body, &gzips); err != nil {
		return nil, err
	}
	return gzips, nil
}

// CreateGzipInput is used as input to the CreateGzip function.
type CreateGzipInput struct {
	// CacheCondition is the name of the cache condition controlling when this configuration applies.
	CacheCondition *string `url:"cache_condition,omitempty"`
	// ContentTypes is a space-separated list of content types to compress.
	ContentTypes *string `url:"content_types,omitempty"`
	// Extensions is a space-separated list of file extensions to compress.
	Extensions *string `url:"extensions,omitempty"`
	// Name is the name of the gzip configuration.
	Name *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// CreateGzip creates a new resource.
func (c *Client) CreateGzip(ctx context.Context, i *CreateGzipInput) (*Gzip, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "gzip")
	resp, err := c.PostForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gzip *Gzip
	if err := DecodeBodyMap(resp.Body, &gzip); err != nil {
		return nil, err
	}
	return gzip, nil
}

// GetGzipInput is used as input to the GetGzip function.
type GetGzipInput struct {
	// Name is the name of the Gzip to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetGzip retrieves the specified resource.
func (c *Client) GetGzip(ctx context.Context, i *GetGzipInput) (*Gzip, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "gzip", i.Name)
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Gzip
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateGzipInput is used as input to the UpdateGzip function.
type UpdateGzipInput struct {
	// CacheCondition is the name of the cache condition controlling when this configuration applies.
	CacheCondition *string `url:"cache_condition,omitempty"`
	// ContentTypes is a space-separated list of content types to compress.
	ContentTypes *string `url:"content_types,omitempty"`
	// Extensions is a space-separated list of file extensions to compress.
	Extensions *string `url:"extensions,omitempty"`
	// Name is the name of the Gzip to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
}

// UpdateGzip updates the specified resource.
func (c *Client) UpdateGzip(ctx context.Context, i *UpdateGzipInput) (*Gzip, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "gzip", i.Name)
	resp, err := c.PutForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Gzip
	if err := DecodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteGzipInput is the input parameter to DeleteGzip.
type DeleteGzipInput struct {
	// Name is the name of the Gzip to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteGzip deletes the specified resource.
func (c *Client) DeleteGzip(ctx context.Context, i *DeleteGzipInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "gzip", i.Name)
	resp, err := c.Delete(ctx, path, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := DecodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
