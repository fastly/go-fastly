package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// BlobStorage represents a blob storage response from the Fastly API.
type BlobStorage struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Path              string     `mapstructure:"path"`
	AccountName       string     `mapstructure:"account_name"`
	Container         string     `mapstructure:"container"`
	SASToken          string     `mapstructure:"sas_token"`
	Period            uint       `mapstructure:"period"`
	TimestampFormat   string     `mapstructure:"timestamp_format"`
	CompressionCodec  string     `mapstructure:"compression_codec"`
	GzipLevel         uint       `mapstructure:"gzip_level"`
	PublicKey         string     `mapstructure:"public_key"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	MessageType       string     `mapstructure:"message_type"`
	Placement         string     `mapstructure:"placement"`
	ResponseCondition string     `mapstructure:"response_condition"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// blobStorageByName is a sortable list of blob storages.
type blobStorageByName []*BlobStorage

// Len, Swap, and Less implement the sortable interface.
func (s blobStorageByName) Len() int      { return len(s) }
func (s blobStorageByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s blobStorageByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListBlobStoragesInput is used as input to the ListBlobStorages function.
type ListBlobStoragesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListBlobStorages returns the list of blob storages for the configuration version.
func (c *Client) ListBlobStorages(i *ListBlobStoragesInput) ([]*BlobStorage, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var as []*BlobStorage
	if err := decodeBodyMap(resp.Body, &as); err != nil {
		return nil, err
	}
	sort.Stable(blobStorageByName(as))
	return as, nil
}

// CreateBlobStorageInput is used as input to the CreateBlobStorage function.
type CreateBlobStorageInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `form:"name,omitempty"`
	Path              string `form:"path,omitempty"`
	AccountName       string `form:"account_name,omitempty"`
	Container         string `form:"container,omitempty"`
	SASToken          string `form:"sas_token,omitempty"`
	Period            uint   `form:"period,omitempty"`
	TimestampFormat   string `form:"timestamp_format,omitempty"`
	CompressionCodec  string `form:"compression_codec,omitempty"`
	GzipLevel         uint   `form:"gzip_level,omitempty"`
	PublicKey         string `form:"public_key,omitempty"`
	Format            string `form:"format,omitempty"`
	FormatVersion     uint   `form:"format_version,omitempty"`
	MessageType       string `form:"message_type,omitempty"`
	Placement         string `form:"placement,omitempty"`
	ResponseCondition string `form:"response_condition,omitempty"`
}

// CreateBlobStorage creates a new Fastly blob storage.
func (c *Client) CreateBlobStorage(i *CreateBlobStorageInput) (*BlobStorage, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var a *BlobStorage
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// GetBlobStorageInput is used as input to the GetBlobStorage function.
type GetBlobStorageInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the blob storage to fetch.
	Name string
}

// GetBlobStorage gets the blob storage configuration with the given parameters.
func (c *Client) GetBlobStorage(i *GetBlobStorageInput) (*BlobStorage, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var a *BlobStorage
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// UpdateBlobStorageInput is used as input to the UpdateBlobStorage function.
type UpdateBlobStorageInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the blob storage to update.
	Name string

	NewName           *string `form:"name,omitempty"`
	Path              *string `form:"path,omitempty"`
	AccountName       *string `form:"account_name,omitempty"`
	Container         *string `form:"container,omitempty"`
	SASToken          *string `form:"sas_token,omitempty"`
	Period            *uint   `form:"period,omitempty"`
	TimestampFormat   *string `form:"timestamp_format,omitempty"`
	CompressionCodec  *string `form:"compression_codec,omitempty"`
	GzipLevel         *uint   `form:"gzip_level,omitempty"`
	PublicKey         *string `form:"public_key,omitempty"`
	Format            *string `form:"format,omitempty"`
	FormatVersion     *uint   `form:"format_version,omitempty"`
	MessageType       *string `form:"message_type,omitempty"`
	Placement         *string `form:"placement,omitempty"`
	ResponseCondition *string `form:"response_condition,omitempty"`
}

// UpdateBlobStorage updates a specific blob storage.
func (c *Client) UpdateBlobStorage(i *UpdateBlobStorageInput) (*BlobStorage, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}

	var a *BlobStorage
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}
	return a, nil
}

// DeleteBlobStorageInput is the input parameter to DeleteBlobStorage.
type DeleteBlobStorageInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the blob storage to delete (required).
	Name string
}

// DeleteBlobStorage deletes the given blob storage version.
func (c *Client) DeleteBlobStorage(i *DeleteBlobStorageInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/azureblob/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
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
