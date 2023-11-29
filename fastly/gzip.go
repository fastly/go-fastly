package fastly

import (
	"fmt"
	"net/url"
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
func (c *Client) ListGzips(i *ListGzipsInput) ([]*Gzip, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/gzip", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gzips []*Gzip
	if err := decodeBodyMap(resp.Body, &gzips); err != nil {
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
func (c *Client) CreateGzip(i *CreateGzipInput) (*Gzip, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/gzip", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var gzip *Gzip
	if err := decodeBodyMap(resp.Body, &gzip); err != nil {
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
func (c *Client) GetGzip(i *GetGzipInput) (*Gzip, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/gzip/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Gzip
	if err := decodeBodyMap(resp.Body, &b); err != nil {
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
func (c *Client) UpdateGzip(i *UpdateGzipInput) (*Gzip, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/gzip/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *Gzip
	if err := decodeBodyMap(resp.Body, &b); err != nil {
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
func (c *Client) DeleteGzip(i *DeleteGzipInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/gzip/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
		return ErrNotOK
	}
	return nil
}
