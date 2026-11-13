package fastly

import (
	"context"
	"strconv"
	"time"
)

// FOS represents a FOS response from the Fastly API.
type FOS struct {
	AccessKey         *string    `mapstructure:"access_key"`
	BucketName        *string    `mapstructure:"bucket_name"`
	CompressionCodec  *string    `mapstructure:"compression_codec"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	FileMaxBytes      *int       `mapstructure:"file_max_bytes"`
	Format            *string    `mapstructure:"format"`
	FormatVersion     *int       `mapstructure:"format_version"`
	GzipLevel         *int       `mapstructure:"gzip_level"`
	MessageType       *string    `mapstructure:"message_type"`
	Name              *string    `mapstructure:"name"`
	Path              *string    `mapstructure:"path"`
	Period            *int       `mapstructure:"period"`
	Placement         *string    `mapstructure:"placement"`
	PublicKey         *string    `mapstructure:"public_key"`
	Redundancy        *string    `mapstructure:"redundancy"`
	Region            *string    `mapstructure:"region"`
	ResponseCondition *string    `mapstructure:"response_condition"`
	SecretKey         *string    `mapstructure:"secret_key"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	TimestampFormat   *string    `mapstructure:"timestamp_format"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// ListFOSesInput is used as input to the ListFOSes function.
type ListFOSesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListFOSes retrieves all resources.
func (c *Client) ListFOSes(ctx context.Context, i *ListFOSesInput) ([]*FOS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "fos")
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var foses []*FOS
	if err := DecodeBodyMap(resp.Body, &foses); err != nil {
		return nil, err
	}
	return foses, nil
}

// CreateFOSInput is used as input to the CreateFOS function.
type CreateFOSInput struct {
	//  AccessKey is the access key for your Fastly Object Storage account.
	AccessKey *string `url:"access_key,omitempty"`
	// BucketName is the name of the Fastly Object Storage bucket.
	BucketName *string `url:"bucket_name,omitempty"`
	// CompressionCodec is the codec used for compressing your logs. Valid values are zstd, snappy, and gzip.
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// FileMaxBytes is the maximum number of bytes for each uploaded file. A value of 0 can be used to indicate there is no limit on the size of uploaded files, otherwise the minimum value is 1048576 bytes (1 MiB.).
	FileMaxBytes *int `url:"file_max_bytes,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `url:"gzip_level,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Name is the name of the FOS to create (required).
	Name *string `url:"name,omitempty"`
	// Path is the path to upload logs to.
	Path *string `url:"path,omitempty"`
	// Period is how frequently log files are finalized so they can be available for reading (in seconds).
	Period *int `url:"period,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// PublicKey is a PGP public key that Fastly will use to encrypt your log files before writing them to disk.
	PublicKey *string `url:"public_key,omitempty"`
	// Redundancy is the redundancy level.
	Redundancy *string `url:"redundancy,omitempty"`
	// Region is the region of the Fastly Object Storage bucket.
	Region *string `url:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is the secret key for your Fastly Object Storage account.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `url:"timestamp_format,omitempty"`
}

// CreateFOS creates a new resource.
func (c *Client) CreateFOS(ctx context.Context, i *CreateFOSInput) (*FOS, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "fos")
	resp, err := c.PostForm(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var fos *FOS
	if err := DecodeBodyMap(resp.Body, &fos); err != nil {
		return nil, err
	}
	return fos, nil
}

// GetFOSInput is used as input to the GetFOS function.
type GetFOSInput struct {
	// Name is the name of the FOS to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetFOS retrieves the specified resource.
func (c *Client) GetFOS(ctx context.Context, i *GetFOSInput) (*FOS, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "fos", i.Name)
	resp, err := c.Get(ctx, path, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var fos *FOS
	if err := DecodeBodyMap(resp.Body, &fos); err != nil {
		return nil, err
	}
	return fos, nil
}

// UpdateFOSInput is used as input to the UpdateFOS function.
type UpdateFOSInput struct {
	//  AccessKey is the access key for your Fastly Object Storage account.
	AccessKey *string `json:"access_key,omitempty"`
	// BucketName is the name of the Fastly Object Storage bucket.
	BucketName *string `json:"bucket_name,omitempty"`
	// CompressionCodec is the codec used for compressing your logs. Valid values are zstd, snappy, and gzip.
	CompressionCodec *string `json:"compression_codec,omitempty"`
	// FileMaxBytes is the maximum number of bytes for each uploaded file. A value of 0 can be used to indicate there is no limit on the size of uploaded files, otherwise the minimum value is 1048576 bytes (1 MiB.).
	FileMaxBytes *int `json:"file_max_bytes,omitempty"`
	// Format is a Fastly log format string.
	Format *string `json:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `json:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `json:"gzip_level,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `json:"message_type,omitempty"`
	// Name is the name of the FOS to update (required).
	Name string `json:"-"`
	// NewName is the new name for the resource.
	NewName *string `json:"name,omitempty"`
	// Path is the path to upload logs to.
	Path *string `json:"path,omitempty"`
	// Period is how frequently log files are finalized so they can be available for reading (in seconds).
	Period *int `json:"period,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed. Use
	// NullValue[string]() to reset the endpoint to automatic placement.
	Placement *Nullable[string] `json:"placement,omitempty"`
	// PublicKey is a PGP public key that Fastly will use to encrypt your log files before writing them to disk.
	PublicKey *string `json:"public_key,omitempty"`
	// Redundancy is the redundancy level.
	Redundancy *string `json:"redundancy,omitempty"`
	// Region is the region of the Fastly Object Storage bucket.
	Region *string `json:"region,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `json:"response_condition,omitempty"`
	// SecretKey is the secret key for your Fastly Object Storage account.
	SecretKey *string `json:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `json:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `json:"-"`
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `json:"timestamp_format,omitempty"`
}

// UpdateFOS updates the specified resource.
func (c *Client) UpdateFOS(ctx context.Context, i *UpdateFOSInput) (*FOS, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "fos", i.Name)
	resp, err := c.PutJSON(ctx, path, i, CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var fos *FOS
	if err := DecodeBodyMap(resp.Body, &fos); err != nil {
		return nil, err
	}
	return fos, nil
}

// DeleteFOSInput is the input parameter to DeleteFOS.
type DeleteFOSInput struct {
	// Name is the name of the FOS to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteFOS deletes the specified resource.
func (c *Client) DeleteFOS(ctx context.Context, i *DeleteFOSInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "fos", i.Name)
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
