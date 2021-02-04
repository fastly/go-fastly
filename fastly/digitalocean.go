package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// DigitalOcean represents a DigitalOcean response from the Fastly API.
type DigitalOcean struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	BucketName        string     `mapstructure:"bucket_name"`
	Domain            string     `mapstructure:"domain"`
	AccessKey         string     `mapstructure:"access_key"`
	SecretKey         string     `mapstructure:"secret_key"`
	Path              string     `mapstructure:"path"`
	Period            uint       `mapstructure:"period"`
	GzipLevel         uint       `mapstructure:"gzip_level"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	ResponseCondition string     `mapstructure:"response_condition"`
	MessageType       string     `mapstructure:"message_type"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	Placement         string     `mapstructure:"placement"`
	PublicKey         string     `mapstructure:"public_key"`
	CompressionCodec  string     `mapstructure:"compression_codec"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// digitaloceansByName is a sortable list of DigitalOceans.
type digitaloceansByName []*DigitalOcean

// Len, Swap, and Less implement the sortable interface.
func (d digitaloceansByName) Len() int      { return len(d) }
func (d digitaloceansByName) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d digitaloceansByName) Less(i, j int) bool {
	return d[i].Name < d[j].Name
}

// ListDigitalOceansInput is used as input to the ListDigitalOceans function.
type ListDigitalOceansInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListDigitalOceans returns the list of DigitalOceans for the configuration version.
func (c *Client) ListDigitalOceans(i *ListDigitalOceansInput) ([]*DigitalOcean, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/digitalocean", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var digitaloceans []*DigitalOcean
	if err := decodeBodyMap(resp.Body, &digitaloceans); err != nil {
		return nil, err
	}
	sort.Stable(digitaloceansByName(digitaloceans))
	return digitaloceans, nil
}

// CreateDigitalOceanInput is used as input to the CreateDigitalOcean function.
type CreateDigitalOceanInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `form:"name,omitempty"`
	BucketName        string `form:"bucket_name,omitempty"`
	Domain            string `form:"domain,omitempty"`
	AccessKey         string `form:"access_key,omitempty"`
	SecretKey         string `form:"secret_key,omitempty"`
	Path              string `form:"path,omitempty"`
	Period            uint   `form:"period,omitempty"`
	GzipLevel         uint   `form:"gzip_level,omitempty"`
	Format            string `form:"format,omitempty"`
	MessageType       string `form:"message_type,omitempty"`
	FormatVersion     uint   `form:"format_version,omitempty"`
	ResponseCondition string `form:"response_condition,omitempty"`
	TimestampFormat   string `form:"timestamp_format,omitempty"`
	Placement         string `form:"placement,omitempty"`
	PublicKey         string `form:"public_key,omitempty"`
	CompressionCodec  string `form:"compression_codec,omitempty"`
}

// CreateDigitalOcean creates a new Fastly DigitalOcean.
func (c *Client) CreateDigitalOcean(i *CreateDigitalOceanInput) (*DigitalOcean, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/digitalocean", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var digitalocean *DigitalOcean
	if err := decodeBodyMap(resp.Body, &digitalocean); err != nil {
		return nil, err
	}
	return digitalocean, nil
}

// GetDigitalOceanInput is used as input to the GetDigitalOcean function.
type GetDigitalOceanInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the DigitalOcean to fetch.
	Name string
}

// GetDigitalOcean gets the DigitalOcean configuration with the given parameters.
func (c *Client) GetDigitalOcean(i *GetDigitalOceanInput) (*DigitalOcean, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/digitalocean/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var digitalocean *DigitalOcean
	if err := decodeBodyMap(resp.Body, &digitalocean); err != nil {
		return nil, err
	}
	return digitalocean, nil
}

// UpdateDigitalOceanInput is used as input to the UpdateDigitalOcean function.
type UpdateDigitalOceanInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the DigitalOcean to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	BucketName        *string `form:"bucket_name,omitempty"`
	Domain            *string `form:"domain,omitempty"`
	AccessKey         *string `form:"access_key,omitempty"`
	SecretKey         *string `form:"secret_key,omitempty"`
	Path              *string `form:"path,omitempty"`
	Period            *uint   `form:"period,omitempty"`
	GzipLevel         *uint   `form:"gzip_level,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
	MessageType       *string `form:"message_type,omitempty"`
	TimestampFormat   *string `form:"timestamp_format,omitempty"`
	Placement         *string `form:"placement,omitempty"`
	PublicKey         *string `form:"public_key,omitempty"`
	CompressionCodec  *string `form:"compression_codec,omitempty"`
}

// UpdateDigitalOcean updates a specific DigitalOcean.
func (c *Client) UpdateDigitalOcean(i *UpdateDigitalOceanInput) (*DigitalOcean, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/digitalocean/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var digitalocean *DigitalOcean
	if err := decodeBodyMap(resp.Body, &digitalocean); err != nil {
		return nil, err
	}
	return digitalocean, nil
}

// DeleteDigitalOceanInput is the input parameter to DeleteDigitalOcean.
type DeleteDigitalOceanInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the DigitalOcean to delete (required).
	Name string
}

// DeleteDigitalOcean deletes the given DigitalOcean version.
func (c *Client) DeleteDigitalOcean(i *DeleteDigitalOceanInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/digitalocean/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrStatusNotOk
	}
	return nil
}
