package fastly

import (
	"strconv"
	"time"
)

// BlobStorage represents a blob storage response from the Fastly API.
type BlobStorage struct {
	AccountName       *string    `mapstructure:"account_name"`
	CompressionCodec  *string    `mapstructure:"compression_codec"`
	Container         *string    `mapstructure:"container"`
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
	ResponseCondition *string    `mapstructure:"response_condition"`
	SASToken          *string    `mapstructure:"sas_token"`
	ServiceID         *string    `mapstructure:"service_id"`
	ServiceVersion    *int       `mapstructure:"version"`
	TimestampFormat   *string    `mapstructure:"timestamp_format"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
}

// ListBlobStoragesInput is used as input to the ListBlobStorages function.
type ListBlobStoragesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListBlobStorages retrieves all resources.
func (c *Client) ListBlobStorages(i *ListBlobStoragesInput) ([]*BlobStorage, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "azureblob")
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var as []*BlobStorage
	if err := decodeBodyMap(resp.Body, &as); err != nil {
		return nil, err
	}
	return as, nil
}

// CreateBlobStorageInput is used as input to the CreateBlobStorage function.
type CreateBlobStorageInput struct {
	// AccountName is the unique Azure Blob Storage namespace in which your data objects are stored.
	AccountName *string `url:"account_name,omitempty"`
	// CompressionCodec is the codec used for compressing your logs (valid values are zstd, snappy, and gzip).
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Container is the name of the Azure Blob Storage container in which to store logs.
	Container *string `url:"container,omitempty"`
	// FileMaxBytes is the maximum number of bytes for each uploaded file. A value of 0 can be used to indicate there is no limit on the size of uploaded files, otherwise the minimum value is 1048576 bytes (1 MiB.).
	FileMaxBytes *int `url:"file_max_bytes,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `url:"gzip_level,omitempty"`
	// MessageType is how the message should be formatted.
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
	// SASToken is the Azure shared access signature providing write access to the blob service objects.
	SASToken *string `url:"sas_token,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `url:"timestamp_format,omitempty"`
}

// CreateBlobStorage creates a new resource.
func (c *Client) CreateBlobStorage(i *CreateBlobStorageInput) (*BlobStorage, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "azureblob")
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a *BlobStorage
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// GetBlobStorageInput is used as input to the GetBlobStorage function.
type GetBlobStorageInput struct {
	// Name is the name of the blob storage to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetBlobStorage retrieves the specified resource.
func (c *Client) GetBlobStorage(i *GetBlobStorageInput) (*BlobStorage, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "azureblob", i.Name)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a *BlobStorage
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// UpdateBlobStorageInput is used as input to the UpdateBlobStorage function.
type UpdateBlobStorageInput struct {
	// AccountName is the unique Azure Blob Storage namespace in which your data objects are stored.
	AccountName *string `url:"account_name,omitempty"`
	// CompressionCodec is the codec used for compressing your logs (valid values are zstd, snappy, and gzip).
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Container is the name of the Azure Blob Storage container in which to store logs.
	Container *string `url:"container,omitempty"`
	// FileMaxBytes is the maximum number of bytes for each uploaded file. A value of 0 can be used to indicate there is no limit on the size of uploaded files, otherwise the minimum value is 1048576 bytes (1 MiB.).
	FileMaxBytes *int `url:"file_max_bytes,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `url:"gzip_level,omitempty"`
	// MessageType is how the message should be formatted.
	MessageType *string `url:"message_type,omitempty"`
	// Name is the name of the blob storage to update (required).
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
	// SASToken is the Azure shared access signature providing write access to the blob service objects.
	SASToken *string `url:"sas_token,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `url:"timestamp_format,omitempty"`
}

// UpdateBlobStorage updates the specified resource.
func (c *Client) UpdateBlobStorage(i *UpdateBlobStorageInput) (*BlobStorage, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "azureblob", i.Name)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var a *BlobStorage
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteBlobStorageInput is the input parameter to DeleteBlobStorage.
type DeleteBlobStorageInput struct {
	// Name is the name of the blob storage to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteBlobStorage deletes the specified resource.
func (c *Client) DeleteBlobStorage(i *DeleteBlobStorageInput) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "logging", "azureblob", i.Name)
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
