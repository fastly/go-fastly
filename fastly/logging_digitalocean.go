package fastly

import (
	"strconv"
	"time"
)

// DigitalOcean represents a DigitalOcean response from the Fastly API.
type DigitalOcean struct {
	AccessKey         *string    `mapstructure:"access_key"`
	BucketName        *string    `mapstructure:"bucket_name"`
	CompressionCodec  *string    `mapstructure:"compression_codec"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Domain            *string    `mapstructure:"domain"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	GzipLevel         *int       `mapstructure:"gzip_level"`
	MessageType       *string    `mapstructure:"message_type"`
	Name              *string    `mapstructure:"name"`
	Path              *string    `mapstructure:"path"`
	Period            *int       `mapstructure:"period"`
	Placement         *string    `mapstructure:"placement"`
	PublicKey         *string    `mapstructure:"public_key"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	SecretKey         *string    `mapstructure:"secret_key"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	TimestampFormat   *string    `mapstructure:"timestamp_format"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// ListDigitalOceansInput is used as input to the ListDigitalOceans function.
type ListDigitalOceansInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListDigitalOceans retrieves all resources.
func (c *Client) ListDigitalOceans(i *ListDigitalOceansInput) ([]*DigitalOcean, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "digitalocean")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var digitaloceans []*DigitalOcean
	if err := decodeBodyMap(resp.Body, &digitaloceans); err != nil {
		return nil, err
	}
	return digitaloceans, nil
}

// CreateDigitalOceanInput is used as input to the CreateDigitalOcean function.
type CreateDigitalOceanInput struct {
	// AccessKey is your DigitalOcean Spaces account access key.
	AccessKey *string `url:"access_key,omitempty"`
	// BucketName is the name of the DigitalOcean Space.
	BucketName *string `url:"bucket_name,omitempty"`
	// CompressionCodec is the codec used for compressing your logs (zstd, snappy, gzip).
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Domain is the domain of the DigitalOcean Spaces endpoint.
	Domain *string `url:"domain,omitempty"`
	// Format is aFastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `url:"gzip_level,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Name is the name for the real-time logging configuration.
	Name *string `url:"name,omitempty"`
	// Path is the path to upload logs to.
	Path *string `url:"path,omitempty"`
	// Period is how frequently log files are finalized so they can be available for reading (in seconds).
	Period *int `url:"period,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// PublicKey is a PGP public key that Fastly will use to encrypt your log files before writing them to disk.
	PublicKey *string `url:"public_key,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is your DigitalOcean Spaces account secret key.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `url:"timestamp_format,omitempty"`
}

// CreateDigitalOcean creates a new resource.
func (c *Client) CreateDigitalOcean(i *CreateDigitalOceanInput) (*DigitalOcean, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "digitalocean")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var digitalocean *DigitalOcean
	if err := decodeBodyMap(resp.Body, &digitalocean); err != nil {
		return nil, err
	}
	return digitalocean, nil
}

// GetDigitalOceanInput is used as input to the GetDigitalOcean function.
type GetDigitalOceanInput struct {
	// Name is the name of the DigitalOcean to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetDigitalOcean retrieves the specified resource.
func (c *Client) GetDigitalOcean(i *GetDigitalOceanInput) (*DigitalOcean, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "digitalocean", i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var digitalocean *DigitalOcean
	if err := decodeBodyMap(resp.Body, &digitalocean); err != nil {
		return nil, err
	}
	return digitalocean, nil
}

// UpdateDigitalOceanInput is used as input to the UpdateDigitalOcean function.
type UpdateDigitalOceanInput struct {
	// AccessKey is your DigitalOcean Spaces account access key.
	AccessKey *string `url:"access_key,omitempty"`
	// BucketName is the name of the DigitalOcean Space.
	BucketName *string `url:"bucket_name,omitempty"`
	// CompressionCodec is the codec used for compressing your logs (zstd, snappy, gzip).
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Domain is the domain of the DigitalOcean Spaces endpoint.
	Domain *string `url:"domain,omitempty"`
	// Format is aFastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `url:"gzip_level,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Name is the name of the DigitalOcean to update (required).
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Path is the path to upload logs to.
	Path *string `url:"path,omitempty"`
	// Period is how frequently log files are finalized so they can be available for reading (in seconds).
	Period *int `url:"period,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// PublicKey is a PGP public key that Fastly will use to encrypt your log files before writing them to disk.
	PublicKey *string `url:"public_key,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is your DigitalOcean Spaces account secret key.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `url:"timestamp_format,omitempty"`
}

// UpdateDigitalOcean updates the specified resource.
func (c *Client) UpdateDigitalOcean(i *UpdateDigitalOceanInput) (*DigitalOcean, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "digitalocean", i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var digitalocean *DigitalOcean
	if err := decodeBodyMap(resp.Body, &digitalocean); err != nil {
		return nil, err
	}
	return digitalocean, nil
}

// DeleteDigitalOceanInput is the input parameter to DeleteDigitalOcean.
type DeleteDigitalOceanInput struct {
	// Name is the name of the DigitalOcean to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteDigitalOcean deletes the specified resource.
func (c *Client) DeleteDigitalOcean(i *DeleteDigitalOceanInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "digitalocean", i.Name)
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
		return ErrStatusNotOk
	}
	return nil
}
