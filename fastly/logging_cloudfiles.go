package fastly

import (
	"strconv"
	"time"
)

// Cloudfiles represents a Cloudfiles response from the Fastly API.
type Cloudfiles struct {
	AccessKey         *string    `mapstructure:"access_key"`
	BucketName        *string    `mapstructure:"bucket_name"`
	CompressionCodec  *string    `mapstructure:"compression_codec"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	GzipLevel         *int       `mapstructure:"gzip_level"`
	MessageType       *string    `mapstructure:"message_type"`
	Name              *string    `mapstructure:"name"`
	Path              *string    `mapstructure:"path"`
	Period            *int       `mapstructure:"period"`
	Placement         *string    `mapstructure:"placement"`
	PublicKey         *string    `mapstructure:"public_key"`
	Region            *string    `mapstructure:"region"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	TimestampFormat   *string    `mapstructure:"timestamp_format"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	User              *string    `mapstructure:"user"`
}

// ListCloudfilesInput is used as input to the ListCloudfiles function.
type ListCloudfilesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListCloudfiles retrieves all resources.
func (c *Client) ListCloudfiles(i *ListCloudfilesInput) ([]*Cloudfiles, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "cloudfiles")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cloudfiles []*Cloudfiles
	if err := decodeBodyMap(resp.Body, &cloudfiles); err != nil {
		return nil, err
	}
	return cloudfiles, nil
}

// CreateCloudfilesInput is used as input to the CreateCloudfiles function.
type CreateCloudfilesInput struct {
	// AccessKey is your Cloud Files account access key.
	AccessKey *string `url:"access_key,omitempty"`
	// BucketName is the name of your Cloud Files container.
	BucketName *string `url:"bucket_name,omitempty"`
	// CompressionCodec is the codec used for compressing your logs (zstd, snappy, gzip).
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Format is a Fastly log format string.
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
	// Region is the region to stream logs to.
	Region *string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `url:"timestamp_format,omitempty"`
	// User is the username for your Cloud Files account.
	User *string `url:"user,omitempty"`
}

// CreateCloudfiles creates a new resource.
func (c *Client) CreateCloudfiles(i *CreateCloudfilesInput) (*Cloudfiles, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "cloudfiles")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cloudfiles *Cloudfiles
	if err := decodeBodyMap(resp.Body, &cloudfiles); err != nil {
		return nil, err
	}
	return cloudfiles, nil
}

// GetCloudfilesInput is used as input to the GetCloudfiles function.
type GetCloudfilesInput struct {
	// Name is the name of the Cloudfiles to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetCloudfiles retrieves the specified resource.
func (c *Client) GetCloudfiles(i *GetCloudfilesInput) (*Cloudfiles, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "cloudfiles", i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cloudfiles *Cloudfiles
	if err := decodeBodyMap(resp.Body, &cloudfiles); err != nil {
		return nil, err
	}
	return cloudfiles, nil
}

// UpdateCloudfilesInput is used as input to the UpdateCloudfiles function.
type UpdateCloudfilesInput struct {
	// AccessKey is your Cloud Files account access key.
	AccessKey *string `url:"access_key,omitempty"`
	// BucketName is the name of your Cloud Files container.
	BucketName *string `url:"bucket_name,omitempty"`
	// CompressionCodec is the codec used for compressing your logs (zstd, snappy, gzip).
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `url:"gzip_level,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Name is the name of the Cloudfiles to update (required).
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
	// Region is the region to stream logs to.
	Region *string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `url:"timestamp_format,omitempty"`
	// User is the username for your Cloud Files account.
	User *string `url:"user,omitempty"`
}

// UpdateCloudfiles updates the specified resource.
func (c *Client) UpdateCloudfiles(i *UpdateCloudfilesInput) (*Cloudfiles, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "cloudfiles", i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cloudfiles *Cloudfiles
	if err := decodeBodyMap(resp.Body, &cloudfiles); err != nil {
		return nil, err
	}
	return cloudfiles, nil
}

// DeleteCloudfilesInput is the input parameter to DeleteCloudfiles.
type DeleteCloudfilesInput struct {
	// Name is the name of the Cloudfiles to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteCloudfiles deletes the specified resource.
func (c *Client) DeleteCloudfiles(i *DeleteCloudfilesInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "cloudfiles", i.Name)
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
