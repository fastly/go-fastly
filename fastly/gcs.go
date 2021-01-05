package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// GCS represents an GCS logging response from the Fastly API.
type GCS struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Bucket            string     `mapstructure:"bucket_name"`
	User              string     `mapstructure:"user"`
	SecretKey         string     `mapstructure:"secret_key"`
	Path              string     `mapstructure:"path"`
	Period            uint       `mapstructure:"period"`
	CompressionCodec  string     `mapstructure:"compression_codec"`
	GzipLevel         uint8      `mapstructure:"gzip_level"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	MessageType       string     `mapstructure:"message_type"`
	ResponseCondition string     `mapstructure:"response_condition"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	Placement         string     `mapstructure:"placement"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// gcsesByName is a sortable list of gcses.
type gcsesByName []*GCS

// Len, Swap, and Less implement the sortable interface.
func (s gcsesByName) Len() int      { return len(s) }
func (s gcsesByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s gcsesByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListGCSsInput is used as input to the ListGCSs function.
type ListGCSsInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListGCSs returns the list of gcses for the configuration version.
func (c *Client) ListGCSs(i *ListGCSsInput) ([]*GCS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/gcs", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var gcses []*GCS
	if err := decodeBodyMap(resp.Body, &gcses); err != nil {
		return nil, err
	}
	sort.Stable(gcsesByName(gcses))
	return gcses, nil
}

// CreateGCSInput is used as input to the CreateGCS function.
type CreateGCSInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `form:"name,omitempty"`
	Bucket            string `form:"bucket_name,omitempty"`
	User              string `form:"user,omitempty"`
	SecretKey         string `form:"secret_key,omitempty"`
	Path              string `form:"path,omitempty"`
	Period            uint   `form:"period,omitempty"`
	FormatVersion     uint   `form:"format_version,omitempty"`
	CompressionCodec  string `form:"compression_codec,omitempty"`
	GzipLevel         uint8  `form:"gzip_level,omitempty"`
	Format            string `form:"format,omitempty"`
	MessageType       string `form:"message_type,omitempty"`
	ResponseCondition string `form:"response_condition,omitempty"`
	TimestampFormat   string `form:"timestamp_format,omitempty"`
	Placement         string `form:"placement,omitempty"`
}

// CreateGCS creates a new Fastly GCS.
func (c *Client) CreateGCS(i *CreateGCSInput) (*GCS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/gcs", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var gcs *GCS
	if err := decodeBodyMap(resp.Body, &gcs); err != nil {
		return nil, err
	}
	return gcs, nil
}

// GetGCSInput is used as input to the GetGCS function.
type GetGCSInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the GCS to fetch.
	Name string
}

// GetGCS gets the GCS configuration with the given parameters.
func (c *Client) GetGCS(i *GetGCSInput) (*GCS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/gcs/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var b *GCS
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateGCSInput is used as input to the UpdateGCS function.
type UpdateGCSInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the GCS to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	Bucket            *string `form:"bucket_name,omitempty"`
	User              *string `form:"user,omitempty"`
	SecretKey         *string `form:"secret_key,omitempty"`
	Path              *string `form:"path,omitempty"`
	Period            *uint   `form:"period,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	CompressionCodec  *string `form:"compression_codec,omitempty"`
	GzipLevel         *uint8  `form:"gzip_level,omitempty"`
	Format            *string `form:"format,omitempty"`
	MessageType       *string `form:"message_type,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	TimestampFormat   *string `form:"timestamp_format,omitempty"`
	Placement         *string `form:"placement,omitempty"`
}

// UpdateGCS updates a specific GCS.
func (c *Client) UpdateGCS(i *UpdateGCSInput) (*GCS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/gcs/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var b *GCS
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// DeleteGCSInput is the input parameter to DeleteGCS.
type DeleteGCSInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the GCS to delete (required).
	Name string
}

// DeleteGCS deletes the given GCS version.
func (c *Client) DeleteGCS(i *DeleteGCSInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/gcs/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
